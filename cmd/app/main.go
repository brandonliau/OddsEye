package main

import (
	"OddsEye/internal/repository"
	"OddsEye/internal/service"
	"OddsEye/pkg/config"
	"OddsEye/pkg/database"
	"OddsEye/pkg/logger"

	_ "modernc.org/sqlite"
)

func main() {
	logger := logger.NewStdLogger(logger.LevelInfo)
	cfg := config.NewServiceConfig("./config/config.yml", logger)
	db := database.NewSqliteDB("./database.db", logger)
	repo := repository.NewOddsRepo()
	service := service.NewOddsService(cfg, repo, db, logger)

	service.Start()
}
