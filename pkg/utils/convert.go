package utils

import (
	"encoding/json"
)

func ToJsonArray(data interface{}) []byte {
	jsn, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	return jsn
}
