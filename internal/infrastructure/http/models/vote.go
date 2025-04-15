package models

import (
	"time"

	"github.com/ArtemSind/food_tinder/internal/domain"
	"github.com/google/uuid"
)

// VoteRequest represents the request body for creating or updating a vote
type VoteRequest struct {
	ProductID string `json:"product_id" validate:"required,uuid4"`
	Score     int    `json:"score" validate:"required,min=1,max=5"`
}

// VoteResponse represents the response body for a vote
type VoteResponse struct {
	ID        uuid.UUID `json:"id"`
	SessionID uuid.UUID `json:"session_id"`
	ProductID uuid.UUID `json:"product_id"`
	Score     int       `json:"score"`
	CreatedAt time.Time `json:"created_at"`
}

// ProductScoreResponse represents a product with its aggregated score in the response
type ProductScoreResponse struct {
	ProductID uuid.UUID `json:"product_id"`
	AvgScore  float64   `json:"avg_score"`
	VoteCount int       `json:"vote_count"`
}

// FromDomain converts a domain vote to an HTTP response
func VoteResponseFromDomain(vote *domain.Vote) *VoteResponse {
	return &VoteResponse{
		ID:        vote.ID,
		SessionID: vote.SessionID,
		ProductID: vote.ProductID,
		Score:     vote.Score,
		CreatedAt: vote.CreatedAt,
	}
}

// ToDomain converts an HTTP request to a domain vote
func (r *VoteRequest) ToDomain(sessionID uuid.UUID) (*domain.Vote, error) {
	productID, err := uuid.Parse(r.ProductID)
	if err != nil {
		return nil, err
	}
	return domain.NewVote(sessionID, productID, r.Score)
}

// VoteListResponse represents a list of votes in the response
type VoteListResponse struct {
	Votes []*VoteResponse `json:"votes"`
	Count int             `json:"count"`
}

// FromDomainList converts a list of domain votes to an HTTP response
func VoteListResponseFromDomain(votes []*domain.Vote) *VoteListResponse {
	result := make([]*VoteResponse, len(votes))
	for i, vote := range votes {
		result[i] = VoteResponseFromDomain(vote)
	}
	return &VoteListResponse{
		Votes: result,
		Count: len(result),
	}
}

// ProductScoreResponseFromDomain converts a domain product score to an HTTP response
func ProductScoreResponseFromDomain(score *domain.ProductScore) *ProductScoreResponse {
	return &ProductScoreResponse{
		ProductID: score.ProductID,
		AvgScore:  score.AvgScore,
		VoteCount: score.VoteCount,
	}
}

// ProductScoreListResponse represents a list of product scores in the response
type ProductScoreListResponse struct {
	Scores []*ProductScoreResponse `json:"scores"`
	Count  int                     `json:"count"`
}

// ProductScoreListResponseFromDomain converts a list of domain product scores to an HTTP response
func ProductScoreListResponseFromDomain(scores []*domain.ProductScore) *ProductScoreListResponse {
	result := make([]*ProductScoreResponse, len(scores))
	for i, score := range scores {
		result[i] = ProductScoreResponseFromDomain(score)
	}
	return &ProductScoreListResponse{
		Scores: result,
		Count:  len(result),
	}
}
