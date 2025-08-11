package config

import (
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port string `mapstructure:"port"`
		Host string `mapstructure:"host"`
	} `mapstructure:"server"`

	Admin struct {
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
	} `mapstructure:"admin"`

	Database struct {
		Path string `mapstructure:"path"`
	} `mapstructure:"database"`
}

// LoadConfig 加载配置文件
func LoadConfig() (*Config, error) {
	configPath := "./data/config.toml"

	// 如果配置文件不存在，创建默认配置
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := createDefaultConfig(configPath); err != nil {
			return nil, err
		}
	}

	viper.SetConfigFile(configPath)
	viper.SetConfigType("toml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// createDefaultConfig 创建默认配置文件
func createDefaultConfig(configPath string) error {
	defaultConfig := `[server]
host = "0.0.0.0"
port = "5080"

[admin]
username = "admin"
password = "xiaoz.org"

[database]
path = "./data/registry.db"
`

	return os.WriteFile(configPath, []byte(defaultConfig), 0644)
}
