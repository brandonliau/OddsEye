package processor

import (
	"time"

	"odds-eye/internal/fetcher"
	"odds-eye/internal/transformer"
	"odds-eye/internal/util"

	"odds-eye/pkg/config"
	"odds-eye/pkg/database"
	"odds-eye/pkg/logger"
)

type fixturesProcessor struct {
	config      *config.FixturesConfig
	fetcher     fetcher.Fetcher[fetcher.FixturesFetchJob]
	transformer transformer.Transformer
	db          database.Database
	logger      logger.Logger
}

func NewFixturesProcessor(config *config.FixturesConfig, db database.Database, logger logger.Logger) *fixturesProcessor {
	return &fixturesProcessor{
		config:      config,
		fetcher:     fetcher.NewFixturesFetcher(config, logger),
		transformer: transformer.NewFixturesTransformer(config, db, logger),
		db:          db,
		logger:      logger,
	}
}

func (p *fixturesProcessor) Process() {
	start := time.Now()

	err := p.db.ExecSQLFile("./pkg/database/migrations/all_fixtures.sql")
	if err != nil {
		p.logger.Fatal("Failed to create fixtures table")
	}

	jobs := make(chan fetcher.FixturesFetchJob, p.config.Options.Workers)
	intermediate := make(chan []byte, p.config.Options.Workers)
	results := make(chan int, p.config.Options.Workers)

	util.LaunchWorkers(p.config.Options.Workers, jobs, intermediate, p.fetcher.Fetch)
	util.LaunchWorkers(p.config.Options.Workers, intermediate, results, p.transformer.Transform)

	var tasks []fetcher.FixturesFetchJob
	for sport, leagues := range p.config.Fetcher.SportLeagues {
		for _, league := range leagues {
			tasks = append(tasks, fetcher.FixturesFetchJob{Sport: sport, League: league})
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
