package general

import (
	"fmt"
	"log"
)

func CheckError(e error, errMsg string) {
	if e != nil {
		fmt.Println(errMsg)
		return
	}
}

func CheckFatalError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
