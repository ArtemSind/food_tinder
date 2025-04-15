package domain_test

import (
	"testing"

	"github.com/ArtemSind/food_tinder/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewVote(t *testing.T) {
	t.Run("valid vote", func(t *testing.T) {
		// Arrange
		sessionID := uuid.New()
		productID := uuid.New()
		score := 4

		// Act
		vote, err := domain.NewVote(sessionID, productID, score)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, vote)
		assert.Equal(t, sessionID, vote.SessionID)
		assert.Equal(t, productID, vote.ProductID)
		assert.Equal(t, score, vote.Score)
		assert.NotEqual(t, uuid.Nil, vote.ID)
		assert.False(t, vote.CreatedAt.IsZero())
	})

	t.Run("score too low", func(t *testing.T) {
		// Arrange
		sessionID := uuid.New()
		productID := uuid.New()
		score := 0 // Invalid score

		// Act
		vote, err := domain.NewVote(sessionID, productID, score)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, domain.ErrInvalidScore, err)
		assert.Nil(t, vote)
	})

	t.Run("score too high", func(t *testing.T) {
		// Arrange
		sessionID := uuid.New()
		productID := uuid.New()
		score := 6 // Invalid score

		// Act
		vote, err := domain.NewVote(sessionID, productID, score)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, domain.ErrInvalidScore, err)
		assert.Nil(t, vote)
	})
}
