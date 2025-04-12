package config

import (
	"fmt"

	"odds-eye/pkg/logger"
)

type OddsConfig struct {
	API struct {
		Token   string `yaml:"token"`
		BaseURL string `yaml:"baseurl"`
		Format  string `yaml:"odds_format"`
	} `yaml:"api"`
	Options struct {
		Workers         int `yaml:"num_workers"`
		FixtureBatch    int `yaml:"fixture_batch_size"`
		SportsbookBatch int `yaml:"sportsbook_batch_size"`
	} `yaml:"options"`
	Fetcher struct {
		SportMarkets map[string][]string `yaml:"sport_markets"`
		Sportsbooks  []string            `yaml:"sportsbooks"`
	} `yaml:"fetcher"`
	Transformer struct {
		Insertion     string            `yaml:"insertion"`
		Normalization map[string]string `yaml:"normalization"`
	} `yaml:"transformer"`
	logger logger.Logger `yaml:"-"`
}

func NewOddsConfig(file string, logger logger.Logger) *OddsConfig {
	cfg := &OddsConfig{
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

func (c *OddsConfig) validate() error {
	if c.API.Token == "" {
		return fmt.Errorf("no token provided")
	}
	if c.API.BaseURL == "" {
		return fmt.Errorf("no base url provided")
	}
	if c.API.Format == "" {
		return fmt.Errorf("no odds format provided")
	}
	if c.Options.Workers <= 0 {
		return fmt.Errorf("max workers must be greater than 0")
	}
	if c.Options.FixtureBatch <= 0 {
		return fmt.Errorf("fixtures batch size must be greater than 0")
	}
	if c.Options.SportsbookBatch <= 0 {
		return fmt.Errorf("sportsbook batch size must be greater than 0")
	}
	if c.Fetcher.SportMarkets == nil {
		return fmt.Errorf("no sport markets provided")
	}
	if c.Fetcher.Sportsbooks == nil {
		return fmt.Errorf("no sportsbooks provided")
	}
	if c.Transformer.Insertion == "" {
		return fmt.Errorf("no insertion query provided")
	}
	return nil
}
