package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Config struct {
	App struct {
		Port   string `yaml:"port"`
		AppURL string `yaml:"app_url"`
	} `yaml:"app"`

	Database struct {
		Type       string `yaml:"type"` // "postgres", "sqlite", or "mysql"
		Host       string `yaml:"host"`
		Port       string `yaml:"port"`
		User       string `yaml:"user"`
		Password   string `yaml:"password"`
		Name       string `yaml:"name"`
		SQLiteFile string `yaml:"sqlite_file"`
	} `yaml:"database"`
}

func GenerateDefaultConfig() error {
	defaultConfig := &Config{
		App: struct {
			Port   string `yaml:"port"`
			AppURL string `yaml:"app_url"`
		}{
			Port:   "8080",
			AppURL: "http://localhost:8080",
		},
		Database: struct {
			Type       string `yaml:"type"`
			Host       string `yaml:"host"`
			Port       string `yaml:"port"`
			User       string `yaml:"user"`
			Password   string `yaml:"password"`
			Name       string `yaml:"name"`
			SQLiteFile string `yaml:"sqlite_file"`
		}{
			Type:       "sqlite",
			SQLiteFile: "yggdrasil_mock.db",
		},
	}

	if _, err := os.Stat("config.yaml"); err == nil {
		log.Println("config.yaml already exists, skipping generation.")
		return nil
	}

	if err := saveConfigToFile("config.yaml", defaultConfig); err != nil {
		log.Printf("Failed to create default config file: %v", err)
		return err
	}

	return nil
}

func LoadConfig() (*Config, error) {
	if _, err := os.Stat("config.yaml"); os.IsNotExist(err) {
		log.Println("Config file not found, creating default config.yaml...")
		defaultConfig := &Config{
			App: struct {
				Port   string `yaml:"port"`
				AppURL string `yaml:"app_url"`
			}{
				Port:   "8080",
				AppURL: "http://localhost:8080",
			},
			Database: struct {
				Type       string `yaml:"type"`
				Host       string `yaml:"host"`
				Port       string `yaml:"port"`
				User       string `yaml:"user"`
				Password   string `yaml:"password"`
				Name       string `yaml:"name"`
				SQLiteFile string `yaml:"sqlite_file"`
			}{
				Type:       "sqlite",
				SQLiteFile: "yggdrasil_mock.db",
			},
		}
		if err := saveConfigToFile("config.yaml", defaultConfig); err != nil {
			log.Printf("Failed to create default config file: %v", err)
			return nil, err
		}
		return defaultConfig, nil
	}

	file, err := os.Open("config.yaml")
	if err != nil {
		log.Printf("Failed to open config file: %v", err)
		return nil, err
	}
	defer file.Close()

	var cfg Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		log.Printf("Failed to decode config file: %v", err)
		return nil, err
	}

	return &cfg, nil
}

func saveConfigToFile(filename string, cfg *Config) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := yaml.NewEncoder(file)
	encoder.SetIndent(2)
	return encoder.Encode(cfg)
}
