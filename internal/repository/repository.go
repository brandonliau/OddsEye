package repository

import (
	"sync"

	"OddsEye/internal/model"
)

type Repository interface {
	Fixtures() map[string][]model.SimpleFixture
	Teams(fixtureID string) (string, string)
	Groupings() []model.Grouping
	GroupedSelections(wg *sync.WaitGroup, jobs chan model.Grouping, results chan model.GroupedSelection)
}
