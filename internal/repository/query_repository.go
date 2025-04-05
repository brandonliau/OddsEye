package repository

import (
	"sync"

	"OddsEye/internal/model"
	"OddsEye/internal/util"

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
		r.logger.Error("Failed to retrieve fixtures: %v", err)
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
		r.logger.Error("Failed to retrieve teams: %v", err)
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

func (r *queryRepository) Groupings() []model.Grouping {
	query := "SELECT DISTINCT id, market, grouping_key FROM all_odds"
	rows, err := r.db.Query(query)
	if err != nil {
		r.logger.Error("Failed to retrieve groupings: %v", err)
	}
	defer rows.Close()

	var groupings []model.Grouping
	var id, market, groupingKey string
	for rows.Next() {
		err := rows.Scan(&id, &market, &groupingKey)
		if err != nil {
			r.logger.Error("Failed to scan row: %v", err)
		}
		groupings = append(groupings, model.Grouping{Id: id, Market: market, GroupingKey: groupingKey})
	}
	return groupings
}

func (r *queryRepository) GroupedSelections(wg *sync.WaitGroup, jobs chan model.Grouping, results chan model.GroupedSelection) {
	defer wg.Done()

	query := "SELECT selection, price FROM all_odds WHERE id = ? AND market = ? AND grouping_key = ?"
	var selection string
	var price float64

	for job := range jobs {
		combos := make(map[string]float64)

		rows, err := r.db.Query(query, job.Id, job.Market, job.GroupingKey)
		if err != nil {
			r.logger.Error("Failed to retrieve grouped selections: %v", err)
		}

		for rows.Next() {
			err := rows.Scan(&selection, &price)
			if err != nil {
				r.logger.Error("Failed to scan row: %v", err)
			}
			if maxv, ok := combos[selection]; !ok || price > maxv {
				combos[selection] = price
			}
		}
		results <- model.GroupedSelection{Group: job, Selections: util.MapCopy(combos)}
		rows.Close()
	}
}
