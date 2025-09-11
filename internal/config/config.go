package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	LogLevel     string `json:"logLevel"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	ExtractAfter string `json:"extractAfter"`
	ApiBaseUrl   string `json:"apiBaseUrl"`
}

func GetWorkDirectory() (string, error) {
	executablePath, _ := os.Executable()
	if strings.Contains(executablePath, "/tmp/") {
		return os.Getwd()
	}

	workDirectory := filepath.Dir(executablePath)
	if strings.Contains(executablePath, "/build/") {
		return strings.Split(executablePath, "/build/")[0], nil
	}
	return workDirectory, nil
}

// LoadJson reads configuration from a JSON file
func LoadJson(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	var cfg Config
	err = json.Unmarshal(data, &cfg)
	return cfg, err
}

func GetConfig() Config {
	// Carrega configuração
	workingDirectory, _ := GetWorkDirectory()
	cfgPath := filepath.Join(workingDirectory, "config.json")
	loadedConfig, err := LoadJson(cfgPath)
	if err != nil {
		panic(err)
	}
	return loadedConfig
}
