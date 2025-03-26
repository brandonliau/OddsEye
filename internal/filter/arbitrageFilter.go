package filter

import (
	"OddsEye/internal/repository"
	"OddsEye/pkg/database"
	"OddsEye/pkg/logger"
)

type arbitrageFilter struct {
	db     database.Database
	repo   repository.Repository
	logger logger.Logger
}

func (f *arbitrageFilter) NewArbitrageFilter(db database.Database, repo repository.Repository, logger logger.Logger) *arbitrageFilter {
	return &arbitrageFilter{
		db:     db,
		repo:   repo,
		logger: logger,
	}
}

func (f *arbitrageFilter) Execute() {
	// fixtures := f.repo.Fixtures()
	return
}
