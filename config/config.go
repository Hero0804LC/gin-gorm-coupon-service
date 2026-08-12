package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config 全局配置结构
type Config struct {
	Server ServerConfig `yaml:"server"`
	Redis  RedisConfig  `yaml:"redis"`
	MySQL  MySQLConfig  `yaml:"mysql"`
}

type ServerConfig struct {
	Addr string `yaml:"addr"`
}

type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type MySQLConfig struct {
	DSN string `yaml:"dsn"`
}

// GlobalConfig 全局配置实例
var GlobalConfig *Config

// Init 加载配置文件
func Init(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open config file failed: %w", err)
	}
	defer file.Close()

	if err := yaml.NewDecoder(file).Decode(&GlobalConfig); err != nil {
		return fmt.Errorf("decode config failed: %w", err)
	}
	return nil
}

// Get 获取全局配置
func Get() *Config {
	return GlobalConfig
}
