package filter

import (
	"math"

	"OddsEye/internal/model"
	"OddsEye/internal/repository"
	"OddsEye/internal/util"

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
	jobs := make(chan model.Grouping, numWorkers)
	results := make(chan model.GroupedSelection, numWorkers)

	groupings := f.repo.Groupings()
	util.LaunchWorkers(numWorkers, jobs, results, f.repo.GroupedSelections)
	util.DistributeJobs(groupings, jobs)

	query := "INSERT INTO expected_value (id, market, selection, grouping_key, price, novig_mult, novig_add, novig_pow, novig_shin, novig_wc) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	for res := range results {
		var prices []float64
		var selections []string
		for k, v := range res.Selections {
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
			f.db.Exec(query, res.Group.Id, res.Group.Market, res.Group.GroupingKey, selections[i], prices[i], m, a, p, s, w)
		}
	}
}
