package repository

import (
	"OddsEye/internal/model"
	"OddsEye/pkg/database"
	"OddsEye/pkg/logger"
)

type queryRepository struct {
	db     database.Database
	logger logger.Logger
}

func NewQueryRepository(db database.Database, logger logger.Logger) *queryRepository {
	return &queryRepository{
		db:     db,
		logger: logger,
	}
}

func (r *queryRepository) Fixtures() map[string][]model.SimpleFixture {
	query := "SELECT id, sport FROM all_fixtures"
	rows, err := r.db.Query(query)
	if err != nil {
		r.logger.Error("Failed to retrieve fixtures from database: %v", err)
	}
	defer rows.Close()

	fixtures := make(map[string][]model.SimpleFixture)
	var id, sport string
	for rows.Next() {
		err := rows.Scan(&id, &sport)
		if err != nil {
			r.logger.Error("Failed to scan row: %v", err)
		}
		fixtures[sport] = append(fixtures[sport], model.SimpleFixture{ID: id, Sport: sport})
	}
	return fixtures
}

func (r *queryRepository) Teams(fixtureID string) (string, string) {
	query := "SELECT home_team, away_team FROM all_fixtures WHERE id = ?"
	rows, err := r.db.Query(query, fixtureID)
	if err != nil {
		r.logger.Error("Failed to retrieve teams from database: %v", err)
	}
	defer rows.Close()

	var home, away string
	for rows.Next() {
		err := rows.Scan(&home, &away)
		if err != nil {
			r.logger.Error("Failed to scan row: %v", err)
		}
	}
	return home, away
}
