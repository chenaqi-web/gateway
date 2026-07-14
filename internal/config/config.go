package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const configPath = "config/config.yaml"

type Config struct {
	Server ServerConfig `yaml:"server"`
	RPC    RPCConfig    `yaml:"rpc"`
	HTTP   HTTPConfig   `yaml:"http"`
}

type ServerConfig struct {
	Addr string `yaml:"addr"`
	Mode string `yaml:"mode"`
}

func Load() (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("read config %s: %w", configPath, err)
	}

	cfg := &Config{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parse config %s: %w", configPath, err)
	}
	return cfg, nil
}
