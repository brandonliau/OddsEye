package config

import (
	"fmt"

	"odds-eye/pkg/logger"
)

type GraderConfig struct {
	API struct {
		Token   string `yaml:"token"`
		BaseURL string `yaml:"baseurl"`
	} `yaml:"api"`
	Options struct {
		ArbThreshold int `yaml:"arb_threshold"`
		EvThreshold int `yaml:"ev_threshold"`
	} `yaml:"options"`
	logger logger.Logger `yaml:"-"`
}

func NewGraderConfig(file string, logger logger.Logger) *GraderConfig {
	cfg := &GraderConfig{
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

func (c *GraderConfig) validate() error {
	if c.API.Token == "" {
		return fmt.Errorf("no token provided")
	}
	if c.API.BaseURL == "" {
		return fmt.Errorf("no base url provided")
	}
	if c.Options.ArbThreshold == 0 {
		return fmt.Errorf("arbitrage threshold must be greater than 0")
	}
	if c.Options.EvThreshold == 0 {
		return fmt.Errorf("expected value threshold must be greater than 0")
	}
	return nil
}
