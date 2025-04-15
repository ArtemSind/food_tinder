package application

import (
	"context"

	"github.com/ArtemSind/food_tinder/internal/domain"
	"github.com/ArtemSind/food_tinder/internal/infrastructure/external/foodji"
	"github.com/google/uuid"
)

type VoteRepository interface {
	Create(ctx context.Context, vote *domain.Vote) error
	Update(ctx context.Context, vote *domain.Vote) error
	GetBySessionID(ctx context.Context, sessionID uuid.UUID) ([]*domain.Vote, error)
	GetAggregatedScores(ctx context.Context) ([]*domain.ProductScore, error)
}

type ProductRepository interface {
	GetProduct(ctx context.Context, id uuid.UUID) (*foodji.Product, error)
	ListProducts(ctx context.Context) ([]uuid.UUID, error)
}

type VoteService struct {
	voteRepo    VoteRepository
	productRepo ProductRepository
}

func NewVoteService(voteRepo VoteRepository, productRepo ProductRepository) *VoteService {
	return &VoteService{
		voteRepo:    voteRepo,
		productRepo: productRepo,
	}
}

func (s *VoteService) CreateOrUpdateVote(ctx context.Context, sessionID, productID uuid.UUID, score int) (*domain.Vote, error) {

	if _, err := s.productRepo.GetProduct(ctx, productID); err != nil {
		return nil, err
	}

	vote, err := domain.NewVote(sessionID, productID, score)
	if err != nil {
		return nil, err
	}

	// Get existing votes for this session/product combo
	votes, err := s.voteRepo.GetBySessionID(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	// Check if vote already exists
	var existingVote *domain.Vote
	for _, v := range votes {
		if v.ProductID == productID {
			existingVote = v
			break
		}
	}

	if existingVote != nil {
		existingVote.Score = score
		if err := s.voteRepo.Update(ctx, existingVote); err != nil {
			return nil, err
		}
		return existingVote, nil
	}

	// Create new vote
	if err := s.voteRepo.Create(ctx, vote); err != nil {
		return nil, err
	}
	return vote, nil
}

func (s *VoteService) GetVotesBySession(ctx context.Context, sessionID uuid.UUID) ([]*domain.Vote, error) {
	return s.voteRepo.GetBySessionID(ctx, sessionID)
}

func (s *VoteService) GetAggregatedScores(ctx context.Context) ([]*domain.ProductScore, error) {
	return s.voteRepo.GetAggregatedScores(ctx)
}
