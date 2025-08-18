package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type Config struct {
	LogLevel     string `json:"logLevel"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	ExtractAfter string `json:"extractAfter"`
	ApiBaseUrl   string `json:"apiBaseUrl"`
}

// GetExecutablePath retorna o caminho apropriado considerando execução com go run ou binário compilado
// GetExecutablePath retorna o caminho do executável (ou o main.go durante go run)
func GetExecutablePath() string {
	execPath, err := os.Executable()
	if err != nil {
		panic(fmt.Errorf("erro ao obter caminho do executável: %w", err))
	}

	execPath, err = filepath.EvalSymlinks(execPath)
	if err != nil {
		panic(fmt.Errorf("erro ao resolver symlink do executável: %w", err))
	}

	// Detecta se está rodando com 'go run`
	if strings.Contains(execPath, os.TempDir()) || strings.Contains(execPath, "go-build") {
		mainFile := getMainGoFilePath()
		if mainFile != "" {
			return mainFile
		}

		// fallback: assume que o primeiro argumento é o.go chamado
		workDir, err := os.Getwd()
		if err != nil {
			panic(fmt.Errorf("erro ao obter diretório de trabalho: %w", err))
		}
		return filepath.Join(workDir, os.Args[0])
	}

	return execPath
}

// GetWorkDirectory retorna a raiz do projeto se for `go run`, ou diretório do executável
func GetWorkDirectory() string {
	execPath := GetExecutablePath()

	// Se for .go rodado via `go run`, extrair raiz do projeto
	if strings.HasSuffix(execPath, ".go") {
		// Procurar por padrão cmd/<alguma-coisa>/main.go
		parts := strings.Split(filepath.ToSlash(execPath), "/")
		for i := len(parts) - 1; i >= 2; i-- {
			if parts[i-2] == "cmd" && parts[i] == "main.go" {
				// raiz = tudo antes do /cmd
				return "/" + filepath.Join(parts[:i-2]...)
			}
		}

		// fallback: diretório do .go
		return filepath.Dir(execPath)
	}

	// Caso normal: binário compilado
	return filepath.Dir(execPath)
}

// getMainGoFilePath tenta encontrar o main.go na pilha de execução
func getMainGoFilePath() string {
	for i := 0; i < 10; i++ {
		_, file, _, ok := runtime.Caller(i)
		if !ok {
			break
		}

		if strings.HasSuffix(file, "main.go") && strings.Contains(file, filepath.FromSlash("cmd/")) {
			return file
		}
	}
	return ""
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
	cfgPath := filepath.Join(GetWorkDirectory(), "config.json")
	loadedConfig, err := LoadJson(cfgPath)
	if err != nil {
		panic(err)
	}
	return loadedConfig
}
