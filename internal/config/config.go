package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const configPath = "conf/config.yaml"

type Config struct {
	Server  ServerConfig  `yaml:"server"`
	Auth    AuthConfig    `yaml:"auth"`
	Email   EmailConfig   `yaml:"Email"`
	RPC     RPCConfig     `yaml:"rpc"`
	HTTP    HTTPConfig    `yaml:"http"`
	Mysql   MySQLConfig   `yaml:"Mysql"`
	Redis   RedisConfig   `yaml:"Redis"`
	Log     LogConfig     `yaml:"Log"`
	LLM     LLMConfig     `yaml:"llm"`
	AiChat  AiChatConfig  `yaml:"ai_chat"`
	Storage StorageConfig `yaml:"storage"`
}

type EmailConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Nickname string `yaml:"nickname"`
	Secret   string `yaml:"secret"`
	From     string `yaml:"from"`
	IsSSL    bool   `yaml:"is_ssl"`
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
