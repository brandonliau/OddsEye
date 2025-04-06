package filter

import (
	"OddsEye/internal/model"
	"OddsEye/internal/repository"
	"OddsEye/internal/util"

	"OddsEye/pkg/database"
	"OddsEye/pkg/logger"
	"OddsEye/pkg/wagermath"
)

type arbitrageFilter struct {
	db     database.Database
	repo   repository.Repository
	logger logger.Logger
}

func NewArbitrageFilter(db database.Database, repo repository.Repository, logger logger.Logger) *arbitrageFilter {
	return &arbitrageFilter{
		db:     db,
		repo:   repo,
		logger: logger,
	}
}

func (f *arbitrageFilter) Execute() {
	jobs := make(chan model.Grouping, numWorkers)
	results := make(chan model.GroupedSelection, numWorkers)

	groupings := f.repo.Groupings()
	util.LaunchWorkers(numWorkers, jobs, results, f.repo.GroupedSelections)
	util.DistributeJobs(groupings, jobs)

	q1 := "INSERT INTO arbitrage (id, market, selection_α, selection_β, price_α, price_β, total_implied_probability, vig) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	q2 := "INSERT INTO arbitrage (id, market, selection_α, selection_β, selection_γ, price_α, price_β, price_γ, total_implied_probability, vig) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	for res := range results {
		if len(res.Selections) == 2 {
			keys := util.Keys(res.Selections)
			x := res.Selections[keys[0]]
			y := res.Selections[keys[1]]
			totalImpliedProb := wagermath.TotalImpliedProbability(x, y)
			vig := totalImpliedProb - 1
			f.db.Exec(q1, res.Group.Id, res.Group.Market, keys[0], keys[1], x, y, totalImpliedProb, vig)
		} else if len(res.Selections) == 3 {
			keys := util.Keys(res.Selections)
			x := res.Selections[keys[0]]
			y := res.Selections[keys[1]]
			z := res.Selections[keys[2]]
			totalImpliedProb := wagermath.TotalImpliedProbability(x, y, z)
			vig := totalImpliedProb - 1
			f.db.Exec(q2, res.Group.Id, res.Group.Market, keys[0], keys[1], keys[2], x, y, z, totalImpliedProb, vig)
		}
	}
}
