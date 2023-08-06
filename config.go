package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	yamlConfig = "sqlc.yaml"
	jsonConfig = "sqlc.json"
)

type PackageConfig struct {
	Name                      string `json:"name" yaml:"name"`
	Path                      string `json:"path" yaml:"path"`
	Engine                    string `json:"engine" yaml:"engine"`
	EmitInterface             bool   `json:"emit_interface" yaml:"emit_interface"`
	EmitResultStructPointers  bool   `json:"emit_result_struct_pointers" yaml:"emit_result_struct_pointers"`
	EmitParamsStructPointers  bool   `json:"emit_params_struct_pointers" yaml:"emit_params_struct_pointers"`
	EmitMethodsWithDBArgument bool   `json:"emit_methods_with_db_argument" yaml:"emit_methods_with_db_argument"`
}

type sqlcConfig struct {
	Packages []PackageConfig `json:"packages" yaml:"packages"`
}

type sqlcConfigVersion struct {
	Version string `json:"version" yaml:"version"`
}

func readConfigV1(name string) (sqlcConfig, error) {
	var cfg sqlcConfig
	f, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	switch name {
	case jsonConfig:
		err = json.NewDecoder(f).Decode(&cfg)
	default:
		err = yaml.NewDecoder(f).Decode(&cfg)
	}
	return cfg, err
}

func readConfig() (sqlcConfig, error) {
	var cfg sqlcConfig
	name, err := configFile()
	if err != nil {
		return cfg, err
	}

	f, err := os.Open(name)
	if err != nil {
		return cfg, err
	}
	defer f.Close()

	var v sqlcConfigVersion
	switch name {
	case jsonConfig:
		err = json.NewDecoder(f).Decode(&v)
		if err != nil {
			return cfg, err
		}
	case yamlConfig:
		err = yaml.NewDecoder(f).Decode(&v)
		if err != nil {
			return cfg, err
		}
	default:
		return cfg, fmt.Errorf("invalid config file %q", name)
	}

	if v.Version == "1" {
		cfg, err = readConfigV1(name)
	} else {
		cfg, err = readConfigV2(name)
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
