// Package config controls authz-checker configuration.
package config

import (
	"encoding/json"
	"io/ioutil"

	defaults "github.com/mcuadros/go-defaults"
)

// AppConfig stores authz-checker configuration
type AppConfig struct {
	Port      int `json:"port,omitempty" yaml:"port,omitempty" ini:"port,omitempty" default:"8080"`
	ToolsPort int `json:"tools_port,omitempty" yaml:"tools_port,omitempty" ini:"tools_port,omitempty" default:"8081"`
}

// ReadAppConfigFromJSON read the Config from a JSON file.
func ReadAppConfigFromJSON(path string) (AppConfig, error) {
	if len(path) == 0 {
		result := &AppConfig{}
		defaults.SetDefaults(result)
		return *result, nil
	}
	//nolint:gosec // not a potential file inclusion.
	// If config doesn't match AppConfig, error will prevent program to go any further
	jsonByte, err := ioutil.ReadFile(path)
	if err != nil {
		return AppConfig{}, err
	}

	var cfg AppConfig

	err = json.Unmarshal(jsonByte, &cfg)
	if err != nil {
		return AppConfig{}, err
	}

	return cfg, nil
}
