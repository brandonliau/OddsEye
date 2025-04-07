package filter

import (
	"math"

	"OddsEye/pkg/database"
	"OddsEye/pkg/logger"
	"OddsEye/pkg/wagermath"
)

type evFilter struct {
	db     database.Database
	logger logger.Logger
}

func NewEvFilter(db database.Database, logger logger.Logger) *evFilter {
	return &evFilter{
		db:     db,
		logger: logger,
	}
}

func (f *evFilter) Filter() {
	fixtureGroups := groupedFixtures(f.db, f.logger)
	selectionGroups := groupedSelections(f.db, f.logger, fixtureGroups)

	query := "INSERT INTO expected_value (id, market, selection, grouping_key, price, novig_mult, novig_add, novig_pow, novig_shin, novig_wc) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	for _, selectionGroup := range selectionGroups {
		var prices []float64
		var selections []string
		for k, v := range selectionGroup.selections {
			prices = append(prices, v)
			selections = append(selections, k)
		}

		if len(selections) != 2 && len(selections) != 3 {
			continue
		}

		mult := wagermath.RemoveVigMultiplicative(prices...)
		add := wagermath.RemoveVigAdditive(prices...)
		pow := wagermath.RemoveVigPower(prices...)
		shin := wagermath.RemoveVigShin(prices...)
		wc := wagermath.RemoveVigWorstCase(mult, add, pow, shin)

		for i := range selections {
			m := math.Round(mult[i]*100) / 100
			a := math.Round(add[i]*100) / 100
			p := math.Round(pow[i]*100) / 100
			s := math.Round(shin[i]*100) / 100
			w := math.Round(wc[i]*100) / 100
			f.db.Exec(query, selectionGroup.grouping.id, selectionGroup.grouping.market, selectionGroup.grouping.groupingKey, selections[i], prices[i], m, a, p, s, w)
		}
	}
}
