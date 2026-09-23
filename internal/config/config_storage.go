package config

type StorageConfig struct {
	Provider  string `yaml:"provider"`   // local | cos | oss
	RootDir   string `yaml:"root_dir"`   //  存储目录
	BaseURL   string `yaml:"base_url"`   // 文件访问域名
	AccessKey string `yaml:"access_key"` // cos / oss
	SecretKey string `yaml:"secret_key"` // cos / oss
	Bucket    string `yaml:"bucket"`
	Region    string `yaml:"region"`   // cos region / oss endpoint
	MaxSize   int64  `yaml:"max_size"` // 单文件最大字节数，0 表示默认 10MB
}
