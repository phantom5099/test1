package config

import (
	"os"
)

// Config 运行时配置
type Config struct {
	Mock        bool
	LLMEndpoint string
	APIKey      string
}

// LoadConfig 从环境变量加载配置，提供合理的默认值
func LoadConfig() *Config {
	cfg := &Config{
		Mock:        true,
		LLMEndpoint: os.Getenv("NEOCODE_LLM_ENDPOINT"),
		APIKey:      os.Getenv("NEOCODE_API_KEY"),
	}
	if v := os.Getenv("NEOCODE_MOCK"); v != "" {
		// 将空值视为 true；非空值启用 mock
		if v == "0" || v == "false" {
			cfg.Mock = false
		} else {
			cfg.Mock = true
		}
	}
	return cfg
}
