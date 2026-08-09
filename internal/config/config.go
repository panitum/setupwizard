package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type App struct {
	Name string `json:"name"`
	Link string `json:"link"`
}

type Config struct {
	Apps []App `json:"apps"`
}

func NewConfig() *Config {
	return &Config{
		Apps: []App{},
	}
}

func LoadConfig() (*Config, error) {
	var cfg *Config

	data, err := os.ReadFile(getFullPath())
	if err != nil {
		return NewConfig(), err
	}

	if err := json.Unmarshal(data, &cfg); err != nil {
		return NewConfig(), err
	}

	return cfg, nil
}

func IsExist() bool {
	if _, err := os.Stat(getFullPath()); err == nil {
		return true
	}

	return false
}

func (cfg *Config) Save() error {
	if err := ensureDir(); err != nil {
		return err
	}

	jsonData, err := json.MarshalIndent(&cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(getFullPath(), jsonData, 0644)
}

func (cfg *Config) FillDummy() {
	*cfg = Config{
		Apps: []App{
			{Name: "app_name_1", Link: "https://example.com/app_name_1.dmg"},
			{Name: "app_name_2", Link: "https://example.com/app_name_2.dmg"},
		},
	}
}

func ensureDir() error {
	dirPath := filepath.Dir(getFullPath())
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return err
	}

	return nil
}

func getFullPath() string {
	return filepath.Join("config", "config.json")
}
