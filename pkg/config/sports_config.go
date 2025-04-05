package config

import (
	"fmt"
	"os"

	"OddsEye/internal/model"

	"OddsEye/pkg/logger"

	"gopkg.in/yaml.v3"
)

type SportsConfig struct {
	Sports map[string]model.Sport `yaml:",inline"`
	logger logger.Logger          `yaml:"-"`
}

func NewSportsConfig(file string, logger logger.Logger) *SportsConfig {
	cfg := &SportsConfig{
		logger: logger,
	}
	err := cfg.load(file)
	if err != nil {
		logger.Fatal("Failed to load config file: %v", err)
	}
	err = cfg.validate()
	if err != nil {
		logger.Fatal("Failed to validate config file: %v", err)
	}
	return cfg
}

func (c *SportsConfig) load(file string) error {
	yamlFile, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("readfile: %v", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		return fmt.Errorf("unmarshal: %v", err)
	}
	return nil
}

func (c *SportsConfig) validate() error {
	if len(c.Sports) == 0 {
		return fmt.Errorf("empty sports")
	}
	return nil
}
