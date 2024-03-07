package core

import (
	"encoding/json"
	"os"
)

func LoadJson[T any](fileName string) (T, error) {
	var data T
	dat, err := os.ReadFile(fileName)
	if err != nil {
		return data, err
	}
	err = json.Unmarshal(dat, &data)
	return data, err
}
