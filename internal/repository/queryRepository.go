package repository

import (
	"OddsEye/pkg/database"
	"OddsEye/pkg/logger"
)

type queryRepository struct {
	db database.Database
	logger logger.Logger
}

func NewQueryRepository(db database.Database, logger logger.Logger) *queryRepository {
	return &queryRepository{
		db: db,
		logger: logger,
	}
}

func (r *queryRepository) Fixtures() []string {
	query := "SELECT id FROM all_fixtures"
	rows, err := r.db.Query(query)
	if err != nil {
		r.logger.Error("Failed to retrieve fixtures from database: %v", err)
	}
	defer rows.Close()
	
	fixtures := make([]string, 0)
	var id string
	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			r.logger.Error("Failed to scan row: %v", err)
		}
		fixtures = append(fixtures, id)
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