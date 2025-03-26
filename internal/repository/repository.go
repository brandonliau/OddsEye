package repository

type Repository interface {
	Fixtures() []string
	Teams(fixtureID string) (string, string)
}
