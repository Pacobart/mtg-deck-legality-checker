package helpers

import (
	"fmt"
)

var DEBUG = false

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func Debug(s string) {
	if DEBUG {
		fmt.Printf("---------\n%s\n---------\n\n", s)
	}
}
