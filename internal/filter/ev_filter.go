package filter

import (
	"OddsEye/internal/repository"

	"OddsEye/pkg/database"
	"OddsEye/pkg/logger"
	"OddsEye/pkg/wagermath"
)

type evFilter struct {
	db     database.Database
	repo   repository.Repository
	logger logger.Logger
}

func NewEvFilter(db database.Database, repo repository.Repository, logger logger.Logger) *evFilter {
	return &evFilter{
		db:     db,
		repo:   repo,
		logger: logger,
	}
}

func (f *evFilter) Execute() {
	query := "SELECT id, market, selection, price, novig_wc FROM fair_odds WHERE price > novig_wc"
	rows, err := f.db.Query(query)
	if err != nil {
		f.logger.Error("Failed to positive ev odds: %v", err)
	}
	defer rows.Close()

	query = "INSERT INTO positive_ev (id, market, selection, price, fair_price, ev) VALUES (?, ?, ?, ?, ?, ?)"
	var id, market, selection string
	var price, wc float64
	for rows.Next() {
		err := rows.Scan(&id, &market, &selection, &price, &wc)
		if err != nil {
			f.logger.Error("Failed to scan row: %v", err)
		}
		ev := (wagermath.ImpliedProbability(wc) * price) - 1
		f.db.Exec(query, id, market, selection, price, wc, ev)
	}
}
