package main

import (
	"fmt"
)

import (
	"github.com/MPRaiden/gator/internal/config"
)

func main() {
	// Set user name on config file
	gatorConfig := config.Read(config.ConfigFileName)
	gatorConfig.SetUser("mpraiden")

	// Read file again after update
	gatorConfig = config.Read(config.ConfigFileName)
	fmt.Println(gatorConfig)

}
