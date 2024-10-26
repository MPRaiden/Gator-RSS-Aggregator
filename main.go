package main

import (
	"fmt"
	"log"
)

import (
	"github.com/MPRaiden/gator/internal/config"
)

func main() {
	// Set user name on config file
	gatorConfig, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}
	gatorConfig.SetUser("mpraiden")

	// Read file again after update
	gatorConfig, err = config.Read()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(gatorConfig.DBURL)
}
