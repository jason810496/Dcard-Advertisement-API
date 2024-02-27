package utils

import (
	"encoding/json"
	"fmt"
	"log"
)

func PrintJson(v interface{}) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(string(b))
}
