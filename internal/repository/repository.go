package repository

import (
	"OddsEye/internal/model"
)

type Repository interface {
	Fixtures() map[string][]model.SimpleFixture
	Teams(fixtureID string) (string, string)
}
