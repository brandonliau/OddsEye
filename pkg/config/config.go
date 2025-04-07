package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config interface {
	validate() error
}

func load[T any](file string, config T) error {
	yamlFile, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("readfile: %v", err)
	}
	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		return fmt.Errorf("unmarshal: %v", err)
	}
	return nil
}
