package main

import (
	"odds-eye/internal/filter"
	"odds-eye/internal/processor"

	"odds-eye/pkg/config"
	"odds-eye/pkg/database"
	"odds-eye/pkg/logger"

	_ "modernc.org/sqlite"
)

func main() {
	logger := logger.NewStdLogger(logger.LevelInfo)
	db := database.NewSqliteDB("./database.db", logger)
	defer db.Close()

	fixturesConfig := config.NewFixturesConfig("./config/fixtures.yml", logger)
	oddsConfig := config.NewOddsConfig("./config/odds.yml", logger)

	fixturesProcessor := processor.NewFixturesProcessor(fixturesConfig, db, logger)
	oddsProcessor := processor.NewOddsProcessor(oddsConfig, db, logger)

	fixturesProcessor.Process()
	oddsProcessor.Process()

	arbFilter := filter.NewArbitrageFilter(db, logger)
	evFilter := filter.NewEvFilter(db, logger)

	arbFilter.Filter()
	evFilter.Filter()
}
