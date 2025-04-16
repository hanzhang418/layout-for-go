package config

import (
	"fmt"
	"github.com/spf13/viper"
	"layout-for-go/pkg/env"
	"path/filepath"
	"sync"
)

// Config 映射 config/*.yaml 中配置的结构体
type Config struct {
	Log struct {
		Level string `mapstructure:"level"`
	} `mapstructure:"log"`

	Default struct {
		Name string `mapstructure:"name"`
		Age  int    `mapstructure:"age"`
	} `mapstructure:"default"`
}

// Init 初始化项目配置
func Init() Config {
	viper.SetConfigFile(filepath.Join("config", fmt.Sprintf("%s.yaml", env.Current())))

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		panic(err)
	}

	return cfg
}

// 初始化全局 Config 变量
var (
	once     sync.Once
	instance Config
)

// Get 单例当前项目配置
func Get() Config {
	once.Do(func() {
		instance = Init()
	})

	return instance
}
