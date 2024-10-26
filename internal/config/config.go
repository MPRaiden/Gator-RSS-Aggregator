package config

import (
	"encoding/json"
	"github.com/MPRaiden/gator/internal/general"
	"os"
)

const ConfigFileName = ".gatorconfig.json"

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func GetConfigFileFullPath(fileName string) (string, error) {
	homeDir, err := os.UserHomeDir()
	general.CheckError(err, "func GetConfigFileFullPath(): error getting user home directory!")

	fileToRead := homeDir + "/" + fileName

	return fileToRead, nil
}

func Read(fileName string) Config {
	fileToRead, err := GetConfigFileFullPath(ConfigFileName)
	general.CheckError(err, "func Read(): error getting config file full path!")

	dat, err := os.ReadFile(fileToRead)
	general.CheckError(err, "func Read(): error reading from file!")

	configStruct := Config{}
	err = json.Unmarshal(dat, &configStruct)
	general.CheckError(err, "func Read(): error unmarshalling struct!")

	return configStruct
}

func (c *Config) SetUser(currentUser string) error {
	c.CurrentUserName = currentUser
	return write(*c)
}

func write(c Config) error {
	fullFilePath, err := GetConfigFileFullPath(ConfigFileName)
	general.CheckError(err, "func write(): error retrieving config file full path!")

	// Marshal the struct to json byte slice
	jsonBytes, err := json.Marshal(c)
	general.CheckError(err, "func write(): error marshalling struct to json!")

	// Write update to config file
	err = os.WriteFile(fullFilePath, jsonBytes, 0644)
	general.CheckError(err, "func write(): error writting json file!")

	return nil
}
