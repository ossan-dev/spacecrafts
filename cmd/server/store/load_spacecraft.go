package store

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"spacecraft/domain"
)

// most of `store/` package (especially tests) can be simplified and or removed by using of `go embed`
// see:
// - https://blog.jetbrains.com/go/2021/06/09/how-to-use-go-embed-in-go-1-16/
// - https://blog.carlmjohnson.net/post/2021/how-to-use-go-embed/
// - https://zetcode.com/golang/embed/

// Nit: should be: LoadSpacecraftFromFile
// Minor: could be LoadSpacecraft(io.Reader) ([]*domain.Spacecraft, error)
// so you can pass anything which can be read: a file or a buffer of []bytes
// (this will save you reading the file content in the tests)
func LoadspacecraftFromFile(filepath string) ([]*domain.Spacecraft, error) {
	// code from line 20 to line 29 can be removed via go embed
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
