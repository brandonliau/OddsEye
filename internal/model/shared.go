package model

type Sport struct {
	Leagues []string `yaml:"leagues"`
	Markets []string `yaml:"markets"`
}

type SimpleFixture struct {
	ID    string
	Sport string
}

type Grouping struct {
	Id          string
	Market      string
	GroupingKey string
}

type GroupedSelection struct {
	Group      Grouping
	Selections map[string]float64
}
