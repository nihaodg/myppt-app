package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config 配置结构
type Config struct {
	Provider  string `json:"provider"`  // "openai" 或 "ollama"
	APIKey    string `json:"api_key"`
	BaseURL   string `json:"base_url"`
	Model     string `json:"model"`
	MaxTokens int    `json:"max_tokens"`
}

var (
	configPath string
	config     *Config
)

// InitConfig 初始化配置
func InitConfig(homeDir string) error {
	configPath = filepath.Join(homeDir, ".oh-my-ppt", "config.json")

	// 初始化默认配置
	config = &Config{
		Provider:  "openai",
		APIKey:    "",
		BaseURL:   "https://api.openai.com/v1",
		Model:     "gpt-4o",
		MaxTokens: 4000,
	}

	// 读取现有配置
	data, err := os.ReadFile(configPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("读取配置文件失败: %w", err)
		}
		// 配置文件不存在，创建默认配置
		return SaveConfig()
	}

	// 解析 JSON
	var loadedCfg Config
	if err := json.Unmarshal(data, &loadedCfg); err != nil {
		// 配置文件损坏，使用默认配置
		return SaveConfig()
	}

	// 更新配置
	*config = loadedCfg

	return nil
}

// SaveConfig 保存配置
func SaveConfig() error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化配置失败: %w", err)
	}

	// 确保目录存在
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建配置目录失败: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %w", err)
	}

	return nil
}

// GetConfig 获取当前配置
func GetConfig() *Config {
	if config == nil {
		config = &Config{
			Provider:  "openai",
			APIKey:    "",
			BaseURL:   "https://api.openai.com/v1",
			Model:     "gpt-4o",
			MaxTokens: 4000,
		}
	}
	return config
}

// UpdateConfig 更新配置项
func UpdateConfig(key, value string) error {
	switch key {
	case "provider":
		if value != "openai" && value != "ollama" {
			return fmt.Errorf("provider 必须是 'openai' 或 'ollama'")
		}
		config.Provider = value
	case "api_key":
		config.APIKey = value
	case "base_url":
		config.BaseURL = value
	case "model":
		config.Model = value
	case "max_tokens":
		var tokens int
		fmt.Sscanf(value, "%d", &tokens)
		config.MaxTokens = tokens
	default:
		return fmt.Errorf("未知配置项: %s", key)
	}
	return SaveConfig()
}

// SetOpenAI 配置 OpenAI
func SetOpenAI(apiKey, baseURL, model string) {
	config.Provider = "openai"
	config.APIKey = apiKey
	config.BaseURL = baseURL
	if model != "" {
		config.Model = model
	}
}

// SetOllama 配置 Ollama
func SetOllama(baseURL, model string) {
	config.Provider = "ollama"
	config.BaseURL = baseURL
	if model != "" {
		config.Model = model
	}
}
