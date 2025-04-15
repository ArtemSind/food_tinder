package domain

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

func NewSession() *Session {
	return &Session{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
	}
}
