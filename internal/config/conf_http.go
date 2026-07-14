package config

type HTTPConfig struct {
	AgentServerAddr string `yaml:"agent_server_addr"`
	RequestTimeout  int    `yaml:"request_timeout"`
}
