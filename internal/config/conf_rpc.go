package config

type RPCConfig struct {
	CoreServerAddr string `yaml:"core_server_addr"`
	RequestTimeout int    `yaml:"request_timeout"`
}
