package config

import (
	"fmt"
	"os"

	"OddsEye/pkg/logger"

	"gopkg.in/yaml.v3"
)

type GeneralConfig struct {
	Token        string            `yaml:"token"`
	Window       int               `yaml:"window"`
	Replacements map[string]string `yaml:"replacements"`
	logger       logger.Logger     `yaml:"-"`
}

func NewGeneralConfig(file string, logger logger.Logger) *GeneralConfig {
	cfg := &GeneralConfig{
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

func (c *GeneralConfig) load(file string) error {
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

func (c *GeneralConfig) validate() error {
	if c.Token == "" {
		return fmt.Errorf("empty token")
	}
	if c.Window == 0 {
		return fmt.Errorf("empty window")
	}
	return nil
}
