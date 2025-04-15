package domain

import (
	"github.com/google/uuid"
)

// ProductScore represents aggregated voting scores for a product
type ProductScore struct {
	ProductID uuid.UUID `json:"product_id"`
	AvgScore  float64   `json:"avg_score"`
	VoteCount int       `json:"vote_count"`
}
