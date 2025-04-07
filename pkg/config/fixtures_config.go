package config

import (
	"fmt"

	"OddsEye/pkg/logger"
)

type FixturesConfig struct {
	API struct {
		Token   string `yaml:"token"`
		BaseURL string `yaml:"baseurl"`
		Window  int    `yaml:"window_size"`
	} `yaml:"api"`
	Options struct {
		Workers int `yaml:"num_workers"`
	} `yaml:"options"`
	Fetcher struct {
		SportLeagues map[string][]string `yaml:"sport_leagues"`
	} `yaml:"fetcher"`
	Transformer struct {
		Insertion string `yaml:"insertion_query"`
	} `yaml:"transformer"`
	logger logger.Logger `yaml:"-"`
}

func NewFixturesConfig(file string, logger logger.Logger) *FixturesConfig {
	cfg := &FixturesConfig{
		logger: logger,
	}
	err := load(file, cfg)
	if err != nil {
		logger.Fatal("Failed to load config file: %v", err)
	}
	err = cfg.validate()
	if err != nil {
		logger.Fatal("Failed to validate config file: %v", err)
	}
	return cfg
}

func (c *FixturesConfig) validate() error {
	if c.API.Token == "" {
		return fmt.Errorf("no token provided")
	}
	if c.API.BaseURL == "" {
		return fmt.Errorf("no base url provided")
	}
	if c.API.Window <= 0 {
		return fmt.Errorf("window size must be greater than 0")
	}
	if c.Options.Workers <= 0 {
		return fmt.Errorf("max workers must be greater than 0")
	}
	if c.Fetcher.SportLeagues == nil {
		return fmt.Errorf("no sport leagues provided")
	}
	if c.Transformer.Insertion == "" {
		return fmt.Errorf("no insertion query provided")
	}
	return nil
}
