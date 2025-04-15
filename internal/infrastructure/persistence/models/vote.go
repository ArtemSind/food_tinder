package models

import (
	"time"

	"github.com/ArtemSind/food_tinder/internal/domain"
	"github.com/google/uuid"
)

// VoteDB represents a vote entity in the database
type VoteDB struct {
	ID        uuid.UUID `db:"id"`
	SessionID uuid.UUID `db:"session_id"`
	ProductID uuid.UUID `db:"product_id"`
	Score     int       `db:"score"`
	CreatedAt time.Time `db:"created_at"`
}

// ToDomain converts a database vote model to a domain vote model
func (v *VoteDB) ToDomain() *domain.Vote {
	return &domain.Vote{
		ID:        v.ID,
		SessionID: v.SessionID,
		ProductID: v.ProductID,
		Score:     v.Score,
		CreatedAt: v.CreatedAt,
	}
}

// FromDomain converts a domain vote model to a database vote model
func VoteFromDomain(vote *domain.Vote) *VoteDB {
	return &VoteDB{
		ID:        vote.ID,
		SessionID: vote.SessionID,
		ProductID: vote.ProductID,
		Score:     vote.Score,
		CreatedAt: vote.CreatedAt,
	}
}

// ToDomainList converts a list of database vote models to domain vote models
func VotesToDomain(votes []*VoteDB) []*domain.Vote {
	result := make([]*domain.Vote, len(votes))
	for i, vote := range votes {
		result[i] = vote.ToDomain()
	}
	return result
}

// FromDomainList converts a list of domain vote models to database vote models
func VotesFromDomain(votes []*domain.Vote) []*VoteDB {
	result := make([]*VoteDB, len(votes))
	for i, vote := range votes {
		result[i] = VoteFromDomain(vote)
	}
	return result
}
