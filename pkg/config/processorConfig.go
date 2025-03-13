package config

import (
	"fmt"
	"os"

	"OddsEye/pkg/logger"

	"gopkg.in/yaml.v3"
)

type ProcessorConfig struct {
	Token       string
	Window      int
	Country     string
	State       string
	Sports      map[string][]string
	Sportsbooks []string
	Sharpbooks  []string
	logger      logger.Logger
}

func NewProcessorConfig(file string, logger logger.Logger) *ProcessorConfig {
	cfg := &ProcessorConfig{
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

func (c *ProcessorConfig) load(file string) error {
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

func (c *ProcessorConfig) validate() error {
	if c.Token == "" {
		return fmt.Errorf("empty token")
	}
	if c.Window == 0 {
		return fmt.Errorf("empty window")
	}
	if len(c.Sports) == 0 {
		return fmt.Errorf("empty sports")
	}
	if len(c.Sportsbooks) == 0 {
		return fmt.Errorf("empty sportsbooks")
	}
	if len(c.Sharpbooks) == 0 {
		return fmt.Errorf("empty sharpbooks")
	}
	return nil
}
