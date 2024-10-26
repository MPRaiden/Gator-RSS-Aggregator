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
	general.CheckError(err)

	fileToRead := homeDir + "/" + fileName

	return fileToRead, nil
}

func Read(fileName string) Config {
	fileToRead, err := GetConfigFileFullPath(ConfigFileName)
	general.CheckError(err)
	dat, err := os.ReadFile(fileToRead)
	general.CheckError(err)

	configStruct := Config{}
	err = json.Unmarshal(dat, &configStruct)
	general.CheckError(err)

	return configStruct
}

func (c Config) SetUser(currentUser string) Config {
	c.CurrentUserName = currentUser
	return c
}
