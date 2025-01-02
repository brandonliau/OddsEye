package service

import (
	"net/http"
	"time"

	"OddsEye/internal/repository"
	"OddsEye/pkg/config"
	"OddsEye/pkg/database"
	"OddsEye/pkg/logger"
	"OddsEye/pkg/retryhttp"
)

type oddsService struct {
	cfg    *config.ServiceConfig
	client *retryhttp.RetryClient
	repo   repository.Repository
	db     database.Database
	logger logger.Logger
}

func NewOddsService(cfg *config.ServiceConfig, repo repository.Repository, db database.Database, logger logger.Logger) *oddsService {
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
	service := &oddsService{
		cfg:    cfg,
		repo:   repo,
		client: client,
		db:     db,
		logger: logger,
	}
	return service
}

func (s *oddsService) migrate() {
	s.db.ExecSQLFile("./pkg/database/migrations/tables.sql")
}

func (s *oddsService) clone() {
	query := "CREATE TABLE all_odds_arbitrage AS SELECT * FROM all_odds"
	s.db.Exec(query)
	query = "CREATE INDEX IF NOT EXISTS id_market_all_odds_arbitrage_idx ON all_odds_arbitrage(id, market)"
	s.db.Exec(query)

	query = "CREATE TABLE all_odds_positive_ev AS SELECT * FROM all_odds"
	s.db.Exec(query)
	query = "CREATE INDEX IF NOT EXISTS id_market_all_odds_positive_ev_idx ON all_odds_positive_ev(id, market)"
	s.db.Exec(query)
}

func (s *oddsService) Start() {
	s.migrate()

	s.SeedFixtures()
	s.SeedFixtureOdds()
	s.SeedFairOdds()
	s.clone()

	s.filter()
	s.ScanArbitrage(100)
}

func (s *oddsService) Stop() {

}
