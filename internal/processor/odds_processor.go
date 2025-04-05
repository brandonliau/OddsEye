package processor

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"OddsEye/internal/model"
	"OddsEye/internal/repository"
	"OddsEye/internal/util"

	"OddsEye/pkg/config"
	"OddsEye/pkg/database"
	"OddsEye/pkg/logger"
	"OddsEye/pkg/retryhttp"
)

type oddsProcessor struct {
	gCfg   *config.GeneralConfig
	sCfg   *config.SportsConfig
	sbCfg  *config.SportsbooksConfig
	client *retryhttp.RetryClient
	db     database.Database
	repo   repository.Repository
	logger logger.Logger
}

type oddsJob struct {
	fixtures    []model.SimpleFixture
	sportsbooks []string
}

func NewOddsProcessor(gCfg *config.GeneralConfig, sCfg *config.SportsConfig, sbCfg *config.SportsbooksConfig, db database.Database, repo repository.Repository, logger logger.Logger) Processor[oddsJob] {
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

	return &oddsProcessor{
		gCfg:   gCfg,
		sCfg:   sCfg,
		sbCfg:  sbCfg,
		client: client,
		db:     db,
		repo:   repo,
		logger: logger,
	}
}

func (p *oddsProcessor) normalizeURL(url string) string {
	for k, v := range p.gCfg.Replacements {
		url = strings.ReplaceAll(url, k, v)
	}
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

func (p *oddsProcessor) fetch(wg *sync.WaitGroup, jobs chan oddsJob, results chan []byte) {
	defer wg.Done()

	baseURL := "https://api.opticodds.com/api/v3/fixtures/odds"

	params := url.Values{}
	params.Add("key", p.gCfg.Token)
	params.Add("odds_format", "decimal")

	for job := range jobs {
		for _, fixture := range job.fixtures {
			params.Add("fixture_id", fixture.ID)
		}
		for _, market := range p.sCfg.Sports[job.fixtures[0].Sport].Markets {
			params.Add("market", market)
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

func (p *oddsProcessor) process(wg *sync.WaitGroup, jobs chan []byte, results chan int) {
	defer wg.Done()

	query := "INSERT INTO all_odds (id, market, selection, sportsbook, price, url, grouping_key) VALUES (?, ?, ?, ?, ?, ?, ?)"
	stmt, _ := p.db.PrepareExec(query)

	for job := range jobs {
		var fixtureOddsWrapper model.FixtureOddsWrapper
		err := json.Unmarshal(job, &fixtureOddsWrapper)
		if err != nil {
			p.logger.Error("Failed to unmarshal fixture odds: %v", err)
		}

		fixtureOdds := fixtureOddsWrapper.Data
		for _, fixtureOdd := range fixtureOdds {
			id := fixtureOdd.ID
			home, away := p.repo.Teams(id)

			for _, odd := range fixtureOdd.Odds {
				_, err := stmt.Exec(
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

	fixtures := p.repo.Fixtures()
	batchedFixtures, batchedSportsbooks := batchJobs(fixtures, p.sbCfg.Sportsbooks["nj"])

	util.LaunchWorkers(numWorkers, jobs, intermediate, p.fetch)
	util.LaunchWorkers(numWorkers, intermediate, results, p.process)

	var tasks []oddsJob
	for _, fixtureBatch := range batchedFixtures {
		for _, sportsbookBatch := range batchedSportsbooks {
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
