package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Http     `yaml:"http"`
	Postgres `yaml:"postgres"`
}

type Http struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type Postgres struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBname   string `yaml:"dbname"`
	SSLmode  string `yaml:"sslmode"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}

	content, err := os.ReadFile(filepath.Join("config", "config.yaml"))
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = yaml.Unmarshal(content, cfg)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal data: %w", err)
	}

	return cfg, nil
}
