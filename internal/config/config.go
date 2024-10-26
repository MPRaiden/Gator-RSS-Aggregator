package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const ConfigFileName = ".gatorconfig.json"

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func check(e error) {
	if e != nil {
		fmt.Println(e)
		return
	}
}

func GetConfigFileFullPath(fileName string) (string, error) {
	homeDir, err := os.UserHomeDir()
	check(err)

	fileToRead := homeDir + "/" + fileName

	return fileToRead, nil
}

func Read(fileName string) Config {
	fileToRead, err := GetConfigFileFullPath(ConfigFileName)
	check(err)
	dat, err := os.ReadFile(fileToRead)
	check(err)

	configStruct := Config{}
	err = json.Unmarshal(dat, &configStruct)
	check(err)

	return configStruct
}

func (c Config) SetUser(currentUser string) Config {
	c.CurrentUserName = currentUser
	return c
}
