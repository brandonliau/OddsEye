package repository

type fixtureTeams struct {
	home string
	away string
}

type oddsRepo struct {
	fixtures []string
	teams    map[string]fixtureTeams
}

func NewOddsRepo() *oddsRepo {
	return &oddsRepo{
		fixtures: make([]string, 0),
		teams:    make(map[string]fixtureTeams),
	}
}

func (repo *oddsRepo) Fixtures() []string {
	return repo.fixtures
}

func (repo *oddsRepo) Teams(fixture string) (string, string) {
	return repo.teams[fixture].home, repo.teams[fixture].away
}

func (repo *oddsRepo) AddFixture(fixture string, home string, away string) {
	repo.fixtures = append(repo.fixtures, fixture)
	repo.teams[fixture] = fixtureTeams{home: home, away: away}
}

func (repo *oddsRepo) RemoveFixture(fixture string) {
	for i, id := range repo.fixtures {
		if id == fixture {
			repo.fixtures = append(repo.fixtures[:i], repo.fixtures[i+1:]...)
			break
		}
	}
	delete(repo.teams, fixture)
}

func (repo *oddsRepo) ClearFixtures() {
	repo.fixtures = nil
}
