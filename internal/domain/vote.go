package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidScore = errors.New("score must be between 1 and 5")
)

type Vote struct {
	ID        uuid.UUID `json:"id"`
	SessionID uuid.UUID `json:"session_id"`
	ProductID uuid.UUID `json:"product_id"`
	Score     int       `json:"score"`
	CreatedAt time.Time `json:"created_at"`
}

func NewVote(sessionID, productID uuid.UUID, score int) (*Vote, error) {
	if score < 1 || score > 5 {
		return nil, ErrInvalidScore
	}

	return &Vote{
		ID:        uuid.New(),
		SessionID: sessionID,
		ProductID: productID,
		Score:     score,
		CreatedAt: time.Now(),
	}, nil
}
