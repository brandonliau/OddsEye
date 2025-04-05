package model

type Sport struct {
	Leagues []string `yaml:"leagues"`
	Markets []string `yaml:"markets"`
}

type SimpleFixture struct {
	ID    string
	Sport string
}
