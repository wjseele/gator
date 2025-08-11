package config

import "os"

type Config struct {
	DB_URL      string  `json:"db_url"`
	CurrentUser *string `json:"current_user_name"`
}

const configFileName = ".gatorconfig.json"

func Read() (Config, error) {
	return Config{}, nil
}

func (cfg Config) SetUser() error {
	return nil
}

func getConfigFilePath() (string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configPath := homedir + configFileName
	return configPath, nil
}

func write(cfg Config) error {
	return nil
}
