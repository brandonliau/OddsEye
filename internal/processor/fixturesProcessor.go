package processor

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"OddsEye/internal/model"
	"OddsEye/internal/util"

	"OddsEye/pkg/config"
	"OddsEye/pkg/database"
	"OddsEye/pkg/logger"
	"OddsEye/pkg/retryhttp"
)

type fixturesProcessor struct {
	cfg    *config.ProcessorConfig
	client *retryhttp.RetryClient
	query  *sql.Stmt
	db     database.Database
	logger logger.Logger
}

type fixtureJob struct {
	sport  string
	league string
}

func NewFixturesProcessor(cfg *config.ProcessorConfig, db database.Database, logger logger.Logger) *fixturesProcessor {
	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout:     90 * time.Second,
	}
	httpClient := &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}
	client := retryhttp.NewRetryClient(httpClient, logger)

	query := "INSERT INTO all_fixtures (id, start_date, home_team, away_team, sport, league) VALUES (?, ?, ?, ?, ?, ?)"
	stmt, _ := db.PrepareExec(query)

	return &fixturesProcessor{
		cfg:    cfg,
		client: client,
		query:  stmt,
		db:     db,
		logger: logger,
	}
}

func (p *fixturesProcessor) fetch(wg *sync.WaitGroup, jobs chan fixtureJob, results chan []byte) {
	defer wg.Done()

	baseURL := "https://api.opticodds.com/api/v3/fixtures/active"
	currentTime := time.Now()
	windowTime := currentTime.Add(time.Duration(p.cfg.Window*24) * time.Hour)

	params := url.Values{}
	params.Set("key", p.cfg.Token)
	params.Set("start_date_after", currentTime.Format(time.RFC3339))
	params.Set("start_date_before", windowTime.Format(time.RFC3339))

	for job := range jobs {
		params.Set("sport", job.sport)
		params.Set("league", job.league)

		data, err := fetchData(baseURL, params, p.client)
		if err != nil {
			p.logger.Error("Failed to fetch fixture data: %v", err)
		}

		var fixturesWrapper models.FixturesWrapper
		err = json.Unmarshal(data, &fixturesWrapper)
		if err != nil {
			p.logger.Error("Failed to unmarshal fixtures: %v", err)
		}

		// handle pagination
		for fixturesWrapper.Page < fixturesWrapper.TotalPages {
			fixturesWrapper.Page += 1
			params.Set("page", strconv.Itoa(fixturesWrapper.Page))

			body, err := fetchData(baseURL, params, p.client)
			if err != nil {
				p.logger.Error("Failed to fetch fixture data: %v", err)
			}

			data = append(data, body...)
			err = json.Unmarshal(body, &fixturesWrapper)
			if err != nil {
				p.logger.Error("Failed to unmarshal fixtures: %v", err)
			}
		}

		results <- data
	}
}

func (p *fixturesProcessor) process(wg *sync.WaitGroup, jobs chan []byte, results chan int) {
	defer wg.Done()

	for job := range jobs {
		var fixturesWrapper models.FixturesWrapper
		err := json.Unmarshal(job, &fixturesWrapper)
		if err != nil {
			p.logger.Error("Failed to unmarshal fixtures: %v", err)
		}
		fixtures := fixturesWrapper.Data

		for _, fixture := range fixtures {
			if !fixture.HasOdds {
				continue
			}
			_, err := p.query.Exec(fixture.ID, fixture.StartDate, fixture.HomeTeam, fixture.AwayTeam, fixture.Sport.ID, fixture.League.ID)
			if err != nil {
				p.logger.Error("Failed to execute fixtures insertion statement: %v", err)
				continue
			}

			results <- 0
		}
	}
}

func (p *fixturesProcessor) Execute() {
	start := time.Now()
	
	jobs := make(chan fixtureJob, numWorkers)
	intermediate := make(chan []byte, numWorkers)
	results := make(chan int, numWorkers)

	util.LaunchWorkers(numWorkers, jobs, intermediate, p.fetch)
	util.LaunchWorkers(numWorkers, intermediate, results, p.process)

	tasks := make([]fixtureJob, 0)
	for sport, leagues := range p.cfg.Sports {
		for _, league := range leagues {
			tasks = append(tasks, fixtureJob{sport, league})
		}
	}

	p.db.Begin()
	util.DistributeJobs(tasks, jobs)
	p.db.Commit()

	var processed int
	for range results {
		processed++
	}
	p.logger.Info("Processed %d fixtures in %v", processed, time.Since(start))	
}
