package filter

import (
	"OddsEye/internal/util"
	"OddsEye/pkg/database"
	"OddsEye/pkg/logger"
)

type Filter interface {
	Filter()
}

type fixtureGrouping struct {
	id          string
	market      string
	groupingKey string
}

func groupedFixtures(db database.Database, logger logger.Logger) []fixtureGrouping {
	query := "SELECT DISTINCT id, market, grouping_key FROM all_odds"
	rows, err := db.Query(query)
	if err != nil {
		logger.Error("Failed to retrieve groupings: %v", err)
	}
	defer rows.Close()

	var groupings []fixtureGrouping
	var id, market, groupingKey string
	for rows.Next() {
		err := rows.Scan(&id, &market, &groupingKey)
		if err != nil {
			logger.Error("Failed to scan row: %v", err)
		}
		groupings = append(groupings, fixtureGrouping{id, market, groupingKey})
	}
	return groupings
}

type selectionGrouping struct {
	grouping   fixtureGrouping
	selections map[string]float64
}

func groupedSelections(db database.Database, logger logger.Logger, fixtureGroupings []fixtureGrouping) []selectionGrouping {
	query := "SELECT selection, price FROM all_odds WHERE id = ? AND market = ? AND grouping_key = ?"

	var groupings []selectionGrouping
	var selection string
	var price float64
	for _, group := range fixtureGroupings {
		combos := make(map[string]float64)

		rows, err := db.Query(query, group.id, group.market, group.groupingKey)
		if err != nil {
			logger.Error("Failed to retrieve grouped selections: %v", err)
		}

		for rows.Next() {
			err := rows.Scan(&selection, &price)
			if err != nil {
				logger.Error("Failed to scan row: %v", err)
			}
			if maxv, ok := combos[selection]; !ok || price > maxv {
				combos[selection] = price
			}
		}

		groupings = append(groupings, selectionGrouping{group, util.MapCopy(combos)})
		rows.Close()
	}
	return groupings
}
