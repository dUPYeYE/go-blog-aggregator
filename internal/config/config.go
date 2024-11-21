package config

import (
	"encoding/json"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DatabaseURL string `json:"db_url"`
	Username    string `json:"current_user_name"`
}

func getConfigFilePath() string {
	homeDir, _ := os.UserHomeDir()
	return homeDir + "/" + configFileName
}

func Read() (Config, error) {
	configPath := getConfigFilePath()

	rawFile, err := os.ReadFile(configPath)
	if err != nil {
		return Config{}, err
	}

	var config Config
	json.Unmarshal(rawFile, &config)
	return config, nil
}

func (c Config) SetUser(username string) error {
	configPath := getConfigFilePath()

	c, err := Read()
	if err != nil {
		return err
	}
	c.Username = username

	configJson, err := json.Marshal(c)
	if err != nil {
		return err
	}
	if err = os.WriteFile(configPath, configJson, 0644); err != nil {
		return err
	}

	return nil
}
