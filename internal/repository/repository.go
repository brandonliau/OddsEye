package repository

type Repository interface {
	Fixtures() []string
	Teams(fixture string) (string, string)
	AddFixture(fixture string, home string, away string)
	RemoveFixture(fixture string)
	ClearFixtures()
}
