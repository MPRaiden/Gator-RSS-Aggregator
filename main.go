package main

import (
	"fmt"
)

import (
	"github.com/MPRaiden/gator/internal/config"
)

func main() {
	gatorConfig := config.Read(".gatorconfig.json")
	fmt.Println(gatorConfig)
}
