package config

import (
	"encoding/json"
	"os"
)

const ConfigFileName = ".gatorconfig.json"

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (Config, error) {
	fileToRead, err := getConfigFileFullPath()
	if err != nil {
		return Config{}, err
	}

	dat, err := os.ReadFile(fileToRead)
	if err != nil {
		return Config{}, err
	}

	configStruct := Config{}
	err = json.Unmarshal(dat, &configStruct)
	if err != nil {
		return Config{}, err
	}

	return configStruct, nil
}

func (c *Config) SetUser(currentUser string) error {
	c.CurrentUserName = currentUser
	return write(*c)
}

func getConfigFileFullPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	fileToRead := homeDir + "/" + ConfigFileName

	return fileToRead, nil
}

func write(c Config) error {
	fullFilePath, err := getConfigFileFullPath()
	if err != nil {
		return err
	}

	jsonBytes, err := json.Marshal(c)
	if err != nil {
		return err
	}

	err = os.WriteFile(fullFilePath, jsonBytes, 0644)
	if err != nil {
		return err
	}

	return nil
}
