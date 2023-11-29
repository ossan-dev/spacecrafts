package store

import (
	"encoding/json"
	"io"
	"os"

	"esdemov8/domain"
)

func LoadSpacecraftsFromFile(filepath string) (res []*domain.Spacecraft, err error) {
	file, err := os.Open(filepath)
	if err != nil {
		return
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &res)
	return
}
