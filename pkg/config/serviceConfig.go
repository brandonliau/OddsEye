package config

import (
	"fmt"
	"os"

	"OddsEye/pkg/logger"

	"gopkg.in/yaml.v3"
)

type ServiceConfig struct {
	Token       string
	Window      int
	Country     string
	State       string
	Sports      map[string][]string
	Sportsbooks []string
	Sharpbooks  []string
	logger      logger.Logger
}

func NewServiceConfig(file string, logger logger.Logger) *ServiceConfig {
	cfg := &ServiceConfig{
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

func (c *ServiceConfig) load(file string) error {
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

func (c *ServiceConfig) validate() error {
	if c.Token == "" {
		return fmt.Errorf("empty token")
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
