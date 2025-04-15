package models

import (
	"time"

	"github.com/ArtemSind/food_tinder/internal/domain"
	"github.com/google/uuid"
)

// SessionResponse represents the response body for a session
type SessionResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

// FromDomain converts a domain session to an HTTP response
func SessionResponseFromDomain(session *domain.Session) *SessionResponse {
	return &SessionResponse{
		ID:        session.ID,
		CreatedAt: session.CreatedAt,
	}
}

// SessionListResponse represents a list of sessions in the response
type SessionListResponse struct {
	Sessions []*SessionResponse `json:"sessions"`
	Count    int                `json:"count"`
}

// FromDomainList converts a list of domain sessions to an HTTP response
func SessionListResponseFromDomain(sessions []*domain.Session) *SessionListResponse {
	result := make([]*SessionResponse, len(sessions))
	for i, session := range sessions {
		result[i] = SessionResponseFromDomain(session)
	}
	return &SessionListResponse{
		Sessions: result,
		Count:    len(result),
	}
}
