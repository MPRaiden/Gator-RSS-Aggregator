package config

import (
	"encoding/json"
	"fmt"
	"os"
)

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

func Read(fileName string) Config {
	homeDir, err := os.UserHomeDir()
	check(err)

	fileToRead := homeDir + "/" + fileName
	dat, err := os.ReadFile(fileToRead)
	check(err)

	configStruct := Config{}
	err = json.Unmarshal(dat, &configStruct)
	check(err)

	return configStruct
}
