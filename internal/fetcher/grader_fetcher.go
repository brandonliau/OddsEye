package fetcher

import (
	"odds-eye/pkg/config"
	"odds-eye/pkg/logger"
	"odds-eye/pkg/retryhttp"
)

type graderFetcher struct {
	config *config.GraderConfig
	client *retryhttp.RetryClient
	logger logger.Logger
}
