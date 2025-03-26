package main

import (
	"OddsEye/internal/processor"
	"OddsEye/internal/repository"
	"OddsEye/pkg/config"
	"OddsEye/pkg/database"
	"OddsEye/pkg/logger"

	_ "modernc.org/sqlite"
)

func main() {
	logger := logger.NewStdLogger(logger.LevelDebug)
	cfg := config.NewProcessorConfig("./config/config.yml", logger)
	db := database.NewSqliteDB("./database.db", logger)
	defer db.Close()

	repo := repository.NewQueryRepository(db, logger)
	fixtureProcessor := processor.NewFixturesProcessor(cfg, db, logger)
	oddsProcessor := processor.NewOddsProcessor(cfg, db, repo, logger)

	err := db.ExecSQLFile("./pkg/database/migrations/tables.sql")
	if err != nil {
		logger.Fatal("Failed to perform database migrations: %v", err)
	}
	
	fixtureProcessor.Execute()
	oddsProcessor.Execute()
}
