package config

type EmailConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Nickname string `yaml:"nickname"`
	Secret   string `yaml:"secret"`
	From     string `yaml:"from"`
	IsSSL    bool   `yaml:"is_ssl"`
}
