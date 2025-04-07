package filter

import (
	"OddsEye/internal/util"

	"OddsEye/pkg/database"
	"OddsEye/pkg/logger"
	"OddsEye/pkg/wagermath"
)

type arbitrageFilter struct {
	db     database.Database
	logger logger.Logger
}

func NewArbitrageFilter(db database.Database, logger logger.Logger) *arbitrageFilter {
	return &arbitrageFilter{
		db:     db,
		logger: logger,
	}
}

func (f *arbitrageFilter) Filter() {
	fixtureGroups := groupedFixtures(f.db, f.logger)
	selectionGroups := groupedSelections(f.db, f.logger, fixtureGroups)

	q1 := "INSERT INTO arbitrage (id, market, grouping_key, selection_α, selection_β, price_α, price_β, total_implied_probability, vig) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"
	q2 := "INSERT INTO arbitrage (id, market, grouping_key, selection_α, selection_β, selection_γ, price_α, price_β, price_γ, total_implied_probability, vig) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	for _, selectionGroup := range selectionGroups {
		if len(selectionGroup.selections) == 2 {
			keys := util.Keys(selectionGroup.selections)
			x := selectionGroup.selections[keys[0]]
			y := selectionGroup.selections[keys[1]]
			totalImpliedProb := wagermath.TotalImpliedProbability(x, y)
			vig := totalImpliedProb - 1
			f.db.Exec(q1, selectionGroup.grouping.id, selectionGroup.grouping.market, selectionGroup.grouping.groupingKey, keys[0], keys[1], x, y, totalImpliedProb, vig)
		} else if len(selectionGroup.selections) == 3 {
			keys := util.Keys(selectionGroup.selections)
			x := selectionGroup.selections[keys[0]]
			y := selectionGroup.selections[keys[1]]
			z := selectionGroup.selections[keys[2]]
			totalImpliedProb := wagermath.TotalImpliedProbability(x, y, z)
			vig := totalImpliedProb - 1
			f.db.Exec(q2, selectionGroup.grouping.id, selectionGroup.grouping.market, selectionGroup.grouping.groupingKey, keys[0], keys[1], keys[2], x, y, z, totalImpliedProb, vig)
		}
	}
}
