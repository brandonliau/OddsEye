package main

import (
	"OddsEye/internal/filter"
	"OddsEye/internal/processor"
	"OddsEye/internal/repository"

	"OddsEye/pkg/config"
	"OddsEye/pkg/database"
	"OddsEye/pkg/logger"

	_ "modernc.org/sqlite"
)

func main() {
	logger := logger.NewStdLogger(logger.LevelInfo)
	gCfg := config.NewGeneralConfig("./config/config.yml", logger)
	sCfg := config.NewSportsConfig("./config/sports.yml", logger)
	sbCfg := config.NewSportsbooksConfig("./config/sportsbooks.yml", logger)
	db := database.NewSqliteDB("./database.db", logger)
	defer db.Close()

	repo := repository.NewQueryRepository(db, logger)
	fixtureProcessor := processor.NewFixturesProcessor(gCfg, sCfg, db, logger)
	oddsProcessor := processor.NewOddsProcessor(gCfg, sCfg, sbCfg, db, repo, logger)

	arbFilter := filter.NewArbitrageFilter(db, repo, logger)
	fairFilter := filter.NewFairFilter(db, repo, logger)

	err := db.ExecSQLFile("./pkg/database/migrations/tables.sql")
	if err != nil {
		logger.Fatal("Failed to perform database migrations: %v", err)
	}

	fixtureProcessor.Execute()
	oddsProcessor.Execute()

	arbFilter.Execute()
	fairFilter.Execute()
}
