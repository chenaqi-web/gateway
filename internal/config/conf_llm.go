package config

type LLMConfig struct {
	Provider  string `yaml:"provider"`
	APIKey    string `yaml:"api_key"`
	APIURL    string `yaml:"api_url"`
	ModelID   string `yaml:"model_id"`
	ModelName string `yaml:"model_name"`
}
