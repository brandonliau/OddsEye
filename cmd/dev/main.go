package main

import (
	"fmt"

	"OddsEye/pkg/config"
	"OddsEye/pkg/logger"
)

func main() {
	logger := logger.NewStdLogger(logger.LevelDebug)
	fCfg := config.NewFixturesConfig("./config/fixtures.yml", logger)
	fmt.Printf("%+v", fCfg)
}
