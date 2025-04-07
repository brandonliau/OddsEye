package main

import (
	"OddsEye/internal/filter"
	"OddsEye/internal/processor"

	"OddsEye/pkg/config"
	"OddsEye/pkg/database"
	"OddsEye/pkg/logger"

	_ "modernc.org/sqlite"
)

func main() {
	logger := logger.NewStdLogger(logger.LevelInfo)
	db := database.NewSqliteDB("./database.db", logger)
	defer db.Close()

	err := db.ExecSQLFile("./pkg/database/migrations/tables.sql")
	if err != nil {
		logger.Fatal("Failed to perform database migrations: %v", err)
	}

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
