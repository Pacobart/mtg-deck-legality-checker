package helpers

import (
	"bytes"
	"encoding/json"
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

func StructToJson(inputStruct interface{}) string {
	var buffer bytes.Buffer
	json.NewEncoder(&buffer).Encode(&inputStruct)
	return buffer.String()
}
