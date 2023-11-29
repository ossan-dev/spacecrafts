package store

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"spacecraft/domain"
)

func LoadspacecraftFromFile(filepath string) ([]*domain.Spacecraft, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("err while opening file: %v", err)
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("err while reading file: %v", err)
	}
	var res []*domain.Spacecraft
	if err = json.Unmarshal(data, &res); err != nil {
		return nil, fmt.Errorf("err while unmarshaling file: %v", err)
	}
	return res, nil
}
