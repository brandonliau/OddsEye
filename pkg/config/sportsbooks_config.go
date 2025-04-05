package config

import (
	"fmt"
	"os"

	"OddsEye/pkg/logger"

	"gopkg.in/yaml.v3"
)

type SportsbooksConfig struct {
	Sportsbooks map[string][]string `yaml:",inline"`
	logger      logger.Logger       `yaml:"-"`
}

func NewSportsbooksConfig(file string, logger logger.Logger) *SportsbooksConfig {
	cfg := &SportsbooksConfig{
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

func (c *SportsbooksConfig) load(file string) error {
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

func (c *SportsbooksConfig) validate() error {
	if len(c.Sportsbooks) == 0 {
		return fmt.Errorf("empty sportsbooks")
	}
	return nil
}
