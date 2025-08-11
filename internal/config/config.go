package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	DB_URL      string `json:"db_url"`
	CurrentUser string `json:"current_user_name"`
}

const configFileName = ".gatorconfig.json"

func Read() (Config, error) {
	path, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}
	config, err := readConfig(path)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}

func (cfg Config) SetUser(userName string) error {
	cfg.CurrentUser = userName
	err := write(cfg)
	if err != nil {
		return err
	}
	return nil
}

func getConfigFilePath() (string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configPath := homedir + "/" + configFileName
	return configPath, nil
}

func readConfig(path string) (Config, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	var config Config
	if err := json.Unmarshal(content, &config); err != nil {
		return Config{}, err
	}
	return config, nil
}

func write(cfg Config) error {
	content, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	path, err := getConfigFilePath()
	if err != nil {
		return err
	}
	err = os.WriteFile(path, content, 0666)
	if err != nil {
		return err
	}
	return nil
}
