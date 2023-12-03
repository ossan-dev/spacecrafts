package store_test

import (
	"encoding/json"
	"os"
	"testing"

	"spacecraft/cmd/server/store"
	"spacecraft/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadspacecraftFromFile(t *testing.T) {
	t.Run("existent file", func(t *testing.T) {
		// Arrange
		spacecraft := []*domain.Spacecraft{
			{
				Uid:        "015",
				Name:       "test spacecraft",
				Registry:   "test registry",
				Status:     "Active",
				DateStatus: "1152",
				SpacecraftClass: &domain.SpacecraftInfo{
					Uid:  "SCM00012",
					Name: "ship",
				},
			},
			{
				Uid:        "022",
				Name:       "test spacecraft 22",
				Registry:   "test registry 22",
				Status:     "Active",
				DateStatus: "0007",
				SpacecraftClass: &domain.SpacecraftInfo{
					Uid:  "SCM0005",
					Name: "barrel",
				},
			},
		}
		// Minor: you have the spacecraft.json file, why don't use it instead of mashalling `spacecraft`
		// write it on a file etc?
		data, err := json.Marshal(&spacecraft)
		require.NoError(t, err)

		// NIT: consider using temp files: https://pkg.go.dev/os#CreateTemp
		filepath := "test_spacecraft.json"
		err = os.WriteFile(filepath, data, os.ModePerm)
		require.NoError(t, err)
		defer func() {
			err = os.Remove(filepath)
			require.NoError(t, err) //nit: not needed to check for this err in test
		}()

		// Act
		res, err := store.LoadspacecraftFromFile(filepath)
		// Assert
		require.NotNil(t, res) // not needed: already have `assert.Equal(t, len(spacecraft), len(res))`
		require.Nil(t, err)
		assert.Equal(t, len(spacecraft), len(res))
	})
	t.Run("non existent file", func(t *testing.T) {
		// Arrange
		// Act
		res, err := store.LoadspacecraftFromFile("unknown_spacecraft.json")
		// Assert
		require.Nil(t, res)
		// Minor: use assert.ErrorContains(t,err, "err while opening file") instead of the below ones
		require.NotNil(t, err)
		assert.Contains(t, err.Error(), "err while opening file")
	})
}
