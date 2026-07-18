package config

type AiChatConfig struct {
	VectorCollection string `yaml:"vector_collection"`
	DefaultTopK      int    `yaml:"default_top_k"`
}
