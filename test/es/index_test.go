package es

import (
	"context"

	"spacecraft/internal/es"

	"github.com/stretchr/testify/require"
)

func (s *EsSuite) TestIndexSpacecraftAsDocuments() {
	s.Run("spacecrafts in ctx", func() {
		// Act
		err := es.IndexSpacecraftAsDocuments(s.Ctx, s.EsClient)
		// Assert
		require.Nil(s.T(), err)
	})

	s.Run("spacecrafts NOT in ctx", func() {
		// Act
		err := es.IndexSpacecraftAsDocuments(context.Background(), s.EsClient)
		// Assert
		require.ErrorContains(s.T(), err, "spacecrafts")
	})
}

func (s *EsSuite) TestIndexSpacecraftAsDocumentsAsync() {
	s.Run("spacecrafts in ctx", func() {
		// Act
		err := es.IndexSpacecraftAsDocumentsAsync(s.Ctx, s.EsClient)
		// Assert
		require.Nil(s.T(), err)
	})

	s.Run("spacecrafts NOT in ctx", func() {
		// Act
		err := es.IndexSpacecraftAsDocumentsAsync(context.Background(), s.EsClient)
		// Assert
		require.ErrorContains(s.T(), err, "spacecrafts")
	})
}

func (s *EsSuite) TestDeleteIndex() {
	s.Run("index name found", func() {
		// Act
		err := es.DeleteIndex(s.EsClient, "spacecrafts")
		// Assert
		require.Nil(s.T(), err)
	})

	s.Run("index name NOT found", func() {
		// Act
		err := es.DeleteIndex(s.EsClient, "unknown")
		// Assert
		require.Nil(s.T(), err)
	})
}
