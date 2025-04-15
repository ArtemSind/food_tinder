package models

import (
	"time"

	"github.com/ArtemSind/food_tinder/internal/domain"
	"github.com/google/uuid"
)

// SessionDB represents a session entity in the database
type SessionDB struct {
	ID        uuid.UUID `db:"id"`
	CreatedAt time.Time `db:"created_at"`
}

// ToDomain converts a database session model to a domain session model
func (s *SessionDB) ToDomain() *domain.Session {
	return &domain.Session{
		ID:        s.ID,
		CreatedAt: s.CreatedAt,
	}
}

// FromDomain converts a domain session model to a database session model
func SessionFromDomain(session *domain.Session) *SessionDB {
	return &SessionDB{
		ID:        session.ID,
		CreatedAt: session.CreatedAt,
	}
}

// ToDomainList converts a list of database session models to domain session models
func SessionsToDomain(sessions []*SessionDB) []*domain.Session {
	result := make([]*domain.Session, len(sessions))
	for i, session := range sessions {
		result[i] = session.ToDomain()
	}
	return result
}

// FromDomainList converts a list of domain session models to database session models
func SessionsFromDomain(sessions []*domain.Session) []*SessionDB {
	result := make([]*SessionDB, len(sessions))
	for i, session := range sessions {
		result[i] = SessionFromDomain(session)
	}
	return result
}
