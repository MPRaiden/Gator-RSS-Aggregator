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
	// Set user name on config file
	gatorConfig := config.Read(config.ConfigFileName)
	gatorConfig.SetUser("Farooo")

	// Marshal the struct to json byte slice
	jsonBytes, err := json.Marshal(gatorConfig)
	general.CheckError(err, "func main(): error while marshalling struct to json!")

	// Retrieve config file full path
	filepath, err := config.GetConfigFileFullPath(config.ConfigFileName)
	general.CheckError(err, "func main(): error while getting config file full path!")

	// Write update to config file
	err = os.WriteFile(filepath, jsonBytes, 0644)
	general.CheckError(err, "func main(): error while writting to file!")

	// Read file again after update
	gatorConfig = config.Read(config.ConfigFileName)
	fmt.Println(gatorConfig)

}
