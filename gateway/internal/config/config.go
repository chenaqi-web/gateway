package config

// Config holds gateway HTTP server and downstream gRPC client settings.
type Config struct {
	Server ServerConfig `yaml:"server"`
	GRPC   GRPCConfig   `yaml:"grpc"`
}

type ServerConfig struct {
	Addr string `yaml:"addr"`
	Mode string `yaml:"mode"` // debug | release | test
}

type GRPCConfig struct {
	CoreServerAddr string `yaml:"core_server_addr"`
}

// Load reads configuration from configs/config.yaml.
func Load() (*Config, error) {
	return &Config{
		Server: ServerConfig{
			Addr: ":8080",
			Mode: "debug",
		},
		GRPC: GRPCConfig{
			CoreServerAddr: "localhost:9090",
		},
	}, nil
}
