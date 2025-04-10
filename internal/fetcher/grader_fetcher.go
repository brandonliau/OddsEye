package fetcher

import (
	"OddsEye/pkg/config"
	"OddsEye/pkg/logger"
	"OddsEye/pkg/retryhttp"
)

type graderFetcher struct {
	config *config.GraderConfig
	client *retryhttp.RetryClient
	logger logger.Logger
}
