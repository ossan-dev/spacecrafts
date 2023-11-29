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
		data, err := json.Marshal(&spacecraft)
		require.NoError(t, err)
		filepath := "test_spacecraft.json"
		err = os.WriteFile(filepath, data, os.ModePerm)
		require.NoError(t, err)
		defer func() {
			err = os.Remove(filepath)
			require.NoError(t, err)
		}()
		// Act
		res, err := store.LoadspacecraftFromFile(filepath)
		// Assert
		require.NotNil(t, res)
		require.Nil(t, err)
		assert.Equal(t, len(spacecraft), len(res))
	})
	t.Run("non existent file", func(t *testing.T) {
		// Arrange
		// Act
		res, err := store.LoadspacecraftFromFile("unknown_spacecraft.json")
		// Assert
		require.Nil(t, res)
		require.NotNil(t, err)
		assert.Contains(t, err.Error(), "err while opening file")
	})
}
