package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"

	"spacecraft/domain"
)

func Printspacecraft(w io.Writer, spacecraft []*domain.Spacecraft) {
	for k, v := range spacecraft {
		fmt.Fprintf(w, "%d\t%v\n", k, v)
	}
}

func WritespacecraftToFile(filepath string, spacecraft []*domain.Spacecraft) error {
	data, err := json.MarshalIndent(spacecraft, "", "\t")
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath, data, fs.ModePerm); err != nil {
		return err
	}
	return nil
}
