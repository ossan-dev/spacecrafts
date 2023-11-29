package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"

	"spacecrafts/domain"
)

func PrintSpacecrafts(w io.Writer, spacecrafts []*domain.Spacecraft) {
	for k, v := range spacecrafts {
		fmt.Fprintf(w, "%d\t%v\n", k, v)
	}
}

func WriteSpacecraftsToFile(filepath string, spacecrafts []*domain.Spacecraft) error {
	data, err := json.MarshalIndent(spacecrafts, "", "\t")
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath, data, fs.ModePerm); err != nil {
		return err
	}
	return nil
}
