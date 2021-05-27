package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const (
	yamlConfig = "sqlc.yaml"
	jsonConfig = "sqlc.json"
)

type sqlcConfig struct {
	Packages []struct {
		Name   string `json:"name" yaml:"name"`
		Path   string `json:"path" yaml:"path"`
		Engine string `json:"engine" yaml:"engine"`
	} `json:"packages" yaml:"packages"`
}

func readConfig() (sqlcConfig, error) {
	var cfg sqlcConfig
	name, err := configFile()
	if err != nil {
		return cfg, err
	}

	f, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	switch name {
	case jsonConfig:
		err = json.NewDecoder(f).Decode(&cfg)
	case yamlConfig:
		err = yaml.NewDecoder(f).Decode(&cfg)
	default:
		return cfg, fmt.Errorf("invalid config file %q", name)
	}
	if err != nil {
		return cfg, err
	}
	for _, pkg := range cfg.Packages {
		if pkg.Name == "" {
			pkg.Name = filepath.Base(pkg.Path)
		}
	}
	return cfg, nil
}

func configFile() (string, error) {
	if f, err := os.Stat(yamlConfig); err == nil && !f.IsDir() {
		return yamlConfig, nil
	}

	if f, err := os.Stat(jsonConfig); err == nil && !f.IsDir() {
		return jsonConfig, nil
	}
	return "", errors.New("no sqlc config files (sqlc.json or sqlc.yaml)")

}
