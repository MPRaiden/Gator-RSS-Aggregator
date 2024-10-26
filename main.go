package main

import (
	"encoding/json"
	"fmt"
	"os"
)

import (
	"github.com/MPRaiden/gator/internal/config"
)

func main() {
	gatorConfig := config.Read(config.ConfigFileName)
	gatorConfig = gatorConfig.SetUser("MPRaiden")

	// Marshal the struct to json byte slice
	jsonBytes, err := json.Marshal(gatorConfig)
	if err != nil {
		fmt.Println("Error while marshaling gatorConfig struct!")
		return
	}

	// Retrieve config file full path
	filepath, err := config.GetConfigFileFullPath(config.ConfigFileName)
	if err != nil {
		fmt.Println("Error while retrieving config file full path!")
		return
	}

	// Write update to config file
	err = os.WriteFile(filepath, jsonBytes, 0644)
	if err != nil {
		fmt.Println("Error while writting to file!")
		return
	}

	// Read file again after update
	gatorConfig = config.Read(config.ConfigFileName)
	fmt.Println(gatorConfig)

}
