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

type fairFilter struct {
	db     database.Database
	repo   repository.Repository
	logger logger.Logger
}

func NewFairFilter(db database.Database, repo repository.Repository, logger logger.Logger) *fairFilter {
	return &fairFilter{
		db:     db,
		repo:   repo,
		logger: logger,
	}
}

func (f *fairFilter) Execute() {
	jobs := make(chan model.Grouping, numWorkers)
	results := make(chan model.GroupedSelection, numWorkers)

	util.LaunchWorkers(numWorkers, jobs, results, f.repo.GroupedSelections)

	groupings := f.repo.Groupings()
	util.DistributeJobs(groupings, jobs)

	query := "INSERT INTO fair_odds (id, market, selection, novig_mult, novig_add, novig_pow, novig_shin, novig_wc, grouping_key) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"
	for res := range results {
		var prices []float64
		var selections []string
		for k, v := range res.Selections {
			selections = append(selections, k)
			prices = append(prices, v)
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
			f.db.Exec(query, res.Group.Id, res.Group.Market, selections[i], m, a, p, s, w, res.Group.GroupingKey)
		}
	}
}
