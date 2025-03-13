package processor

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"OddsEye/internal/model"
	"OddsEye/internal/util"

	"OddsEye/pkg/config"
	"OddsEye/pkg/database"
	"OddsEye/pkg/logger"
	"OddsEye/pkg/retryhttp"
)

type oddsProcessor struct {
	cfg    *config.ProcessorConfig
	client *retryhttp.RetryClient
	query  *sql.Stmt
	db     database.Database
	logger logger.Logger
}

type oddsJob struct {
	fixtureIDs  []string
	sportsbooks []string
}

func NewOddsProcessor(cfg *config.ProcessorConfig, db database.Database, logger logger.Logger) *oddsProcessor {
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

	query := "INSERT INTO all_odds (id, market, selection, sportsbook, price, url, grouping_key) VALUES (?, ?, ?, ?, ?, ?, ?)"
	stmt, _ := db.PrepareExec(query)

	return &oddsProcessor{
		cfg:    cfg,
		client: client,
		query:  stmt,
		db:     db,
		logger: logger,
	}
}

func (p *oddsProcessor) normalizeURL(url string) string {
	url = strings.ReplaceAll(url, "<COUNTRY>", p.cfg.Country)
	url = strings.ReplaceAll(url, "<STATE>", p.cfg.State)
	url = strings.ReplaceAll(url, "blue_book", "fanduel")
	return url
}

func (p *oddsProcessor) groupingKey(normalizedSelection, selection, selectionLine, market, home, away string, points *float64) string {
	groupingKey := "default"
	if normalizedSelection != "" && selectionLine != "" && points != nil {
		groupingKey = fmt.Sprintf("%s:%.1f", normalizedSelection, math.Abs(*points))
	} else if normalizedSelection != "" && selectionLine != "" && (!strings.Contains(market, "correct_score") && !strings.Contains(market, "set_betting")) {
		groupingKey = normalizedSelection
	} else if points != nil {
		if selection == home {
			groupingKey = fmt.Sprintf("default:%.2f", *points)
		} else if selection == away {
			groupingKey = fmt.Sprintf("default:%.2f", -*points)
		} else if selection == "Draw" {
			groupingKey = fmt.Sprintf("default:%.2f", -*points)
		} else {
			groupingKey = fmt.Sprintf("default:%.2f", *points)
		}
	}
	return groupingKey
}

func (p *oddsProcessor) getFixtures() []string {
	query := "SELECT id FROM all_fixtures"
	rows, _ := p.db.Query(query)
	defer rows.Close()

	var id string
	fixtures := make([]string, 0)
	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			p.logger.Error("Failed to scan row: %v", err)
		}
		fixtures = append(fixtures, id)
	}

	return fixtures
}

func (p *oddsProcessor) getTeams(id string) (string, string) {
	query := "SELECT home_team, away_team FROM all_fixtures WHERE id = ?"
	row, _ := p.db.QueryRow(query, id)

	var home, away string
	err := row.Scan(&home, &away)
	if err != nil {
		p.logger.Error("Failed to scan row: %v", err)
	}

	return home, away
}

func (p *oddsProcessor) fetch(wg *sync.WaitGroup, jobs chan oddsJob, results chan []byte) {
	defer wg.Done()

	baseURL := "https://api.opticodds.com/api/v3/fixtures/odds"

	params := url.Values{}
	params.Add("key", p.cfg.Token)
	params.Add("odds_format", "decimal")

	for job := range jobs {
		for _, id := range job.fixtureIDs {
			params.Add("fixture_id", id)
		}
		for _, sportsbook := range job.sportsbooks {
			params.Add("sportsbook", sportsbook)
		}

		body, err := fetchData(baseURL, params, p.client)
		if err != nil {
			p.logger.Error("Failed to fetch odds data: %v", err)
		}

		results <- body
	}
}

// func (p *oddsProcessor) process(data []byte) error {
func (p *oddsProcessor) process(wg *sync.WaitGroup, jobs chan []byte, results chan int) {
	defer wg.Done()

	for job := range jobs {
		var fixtureOddsWrapper models.FixtureOddsWrapper
		err := json.Unmarshal(job, &fixtureOddsWrapper)
		if err != nil {
			p.logger.Error("Failed to unmarshal fixture odds: %v", err)
		}

		fixtureOdds := fixtureOddsWrapper.Data
		for _, fixtureOdd := range fixtureOdds {
			id := fixtureOdd.ID
			home, away := p.getTeams(id)

			for _, odd := range fixtureOdd.Odds {
				_, err := p.query.Exec(
					id,
					odd.Market,
					strings.Split(odd.ID, ":")[3],
					odd.Sportsbook,
					odd.Price,
					p.normalizeURL(odd.DeepLink.Desktop),
					p.groupingKey(odd.NormalizedSelection, odd.Selection, odd.SelectionLine, odd.Market, home, away, odd.Points),
				)
				if err != nil {
					p.logger.Error("Failed to execute odds insertion statement: %v", err)
					continue
				}
				results <- 0
			}
		}
	}
}

func (p *oddsProcessor) Execute() {
	start := time.Now()

	jobs := make(chan oddsJob, numWorkers)
	intermediate := make(chan []byte, numWorkers)
	results := make(chan int, numWorkers)

	data := p.getFixtures()
	fixtures, sportsbooks := batchJobs(data, p.cfg.Sportsbooks)

	util.LaunchWorkers(numWorkers, jobs, intermediate, p.fetch)
	util.LaunchWorkers(numWorkers, intermediate, results, p.process)

	tasks := make([]oddsJob, 0)
	for _, fixtureBatch := range fixtures {
		for _, sportsbookBatch := range sportsbooks {
			tasks = append(tasks, oddsJob{fixtureBatch, sportsbookBatch})
		}
	}

	p.db.Begin()
	util.DistributeJobs(tasks, jobs)
	p.db.Commit()

	var processed int
	for range results {
		processed++
	}
	p.logger.Info("Processed %d odds in %v", processed, time.Since(start))
}
