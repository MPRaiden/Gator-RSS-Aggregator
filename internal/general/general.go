package general

import (
	"fmt"
	"log"
)

func CheckError(e error, errMsg string) {
	if e != nil {
		fmt.Println(errMsg)
		log.Fatal(e)
	}
}
