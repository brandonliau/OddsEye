package processor

import (
	"slices"
	"time"

	"OddsEye/internal/fetcher"
	"OddsEye/internal/transformer"
	"OddsEye/internal/util"

	"OddsEye/pkg/config"
	"OddsEye/pkg/database"
	"OddsEye/pkg/logger"
)

type oddsProcessor struct {
	config      *config.OddsConfig
	fetcher     fetcher.Fetcher[fetcher.OddsFetchJob]
	transformer transformer.Transformer
	db          database.Database
	logger      logger.Logger
}

func NewOddsProcessor(config *config.OddsConfig, db database.Database, logger logger.Logger) *oddsProcessor {
	return &oddsProcessor{
		config:      config,
		fetcher:     fetcher.NewOddsFetcher(config, logger),
		transformer: transformer.NewOddsTransformer(config, db, logger),
		db:          db,
		logger:      logger,
	}
}

func (p *oddsProcessor) fixtures() map[string][]struct {
	id    string
	sport string
} {
	query := "SELECT id, sport FROM all_fixtures"
	rows, err := p.db.Query(query)
	if err != nil {
		p.logger.Error("Failed to retrieve fixtures: %v", err)
	}
	defer rows.Close()

	fixtures := make(map[string][]struct {
		id    string
		sport string
	})
	var id, sport string
	for rows.Next() {
		err := rows.Scan(&id, &sport)
		if err != nil {
			p.logger.Error("Failed to scan row: %v", err)
		}
		fixtures[sport] = append(fixtures[sport], struct {
			id    string
			sport string
		}{id: id, sport: sport})
	}
	return fixtures
}

func (p *oddsProcessor) batchJobs(mapped map[string][]struct {
	id    string
	sport string
}, sportsbooks []string) ([][]string, [][]string, [][]string) {
	// create fixture batches
	var batchFixtures, batchMarkets [][]string
	for sport, fixtures := range mapped {
		for chunk := range slices.Chunk(fixtures, p.config.Options.FixtureBatch) {
			var temp []string
			for _, fixture := range chunk {
				temp = append(temp, fixture.id)
			}
			batchFixtures = append(batchFixtures, temp)
			batchMarkets = append(batchMarkets, p.config.Fetcher.SportMarkets[sport])
		}
	}

	// create sportsbook batches
	var batchSportsbooks [][]string
	for c := range slices.Chunk(sportsbooks, p.config.Options.SportsbookBatch) {
		batchSportsbooks = append(batchSportsbooks, c)
	}

	return batchFixtures, batchMarkets, batchSportsbooks
}

func (p *oddsProcessor) Process() {
	start := time.Now()

	err := p.db.ExecSQLFile("./pkg/database/migrations/all_odds.sql")
	if err != nil {
		p.logger.Fatal("Failed to create odds table")
	}

	jobs := make(chan fetcher.OddsFetchJob, p.config.Options.Workers)
	intermediate := make(chan []byte, p.config.Options.Workers)
	results := make(chan int, p.config.Options.Workers)

	fixtures := p.fixtures()
	batchedFixtures, batchedMarkets, batchedSportsbooks := p.batchJobs(fixtures, p.config.Fetcher.Sportsbooks)

	util.LaunchWorkers(p.config.Options.Workers, jobs, intermediate, p.fetcher.Fetch)
	util.LaunchWorkers(p.config.Options.Workers, intermediate, results, p.transformer.Transform)

	var tasks []fetcher.OddsFetchJob
	for i, fixtureBatch := range batchedFixtures {
		for _, sportsbookBatch := range batchedSportsbooks {
			tasks = append(tasks, fetcher.OddsFetchJob{Fixtures: fixtureBatch, Markets: batchedMarkets[i], Sportsbooks: sportsbookBatch})
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
