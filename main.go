package main

import (
	"encoding/json"
	"fmt"
	"os"
)

import (
	"github.com/MPRaiden/gator/internal/config"
	"github.com/MPRaiden/gator/internal/general"
)

func main() {
	gatorConfig := config.Read(config.ConfigFileName)
	gatorConfig = gatorConfig.SetUser("MPRaiden")

	// Marshal the struct to json byte slice
	jsonBytes, err := json.Marshal(gatorConfig)
	general.CheckError(err)

	// Retrieve config file full path
	filepath, err := config.GetConfigFileFullPath(config.ConfigFileName)
	general.CheckError(err)

	// Write update to config file
	err = os.WriteFile(filepath, jsonBytes, 0644)
	general.CheckError(err)

	// Read file again after update
	gatorConfig = config.Read(config.ConfigFileName)
	fmt.Println(gatorConfig)

}
