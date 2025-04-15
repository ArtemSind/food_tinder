package application_test

import (
	"context"
	"testing"

	"github.com/ArtemSind/food_tinder/internal/application"
	"github.com/ArtemSind/food_tinder/internal/domain"
	"github.com/ArtemSind/food_tinder/internal/infrastructure/external/foodji"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Mock VoteRepository
type MockVoteRepository struct {
	mock.Mock
}

func (m *MockVoteRepository) Create(ctx context.Context, vote *domain.Vote) error {
	args := m.Called(ctx, vote)
	return args.Error(0)
}

func (m *MockVoteRepository) Update(ctx context.Context, vote *domain.Vote) error {
	args := m.Called(ctx, vote)
	return args.Error(0)
}

func (m *MockVoteRepository) GetBySessionID(ctx context.Context, sessionID uuid.UUID) ([]*domain.Vote, error) {
	args := m.Called(ctx, sessionID)
	return args.Get(0).([]*domain.Vote), args.Error(1)
}

func (m *MockVoteRepository) GetAggregatedScores(ctx context.Context) ([]*domain.ProductScore, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*domain.ProductScore), args.Error(1)
}

type MockProductRepository struct {
	mock.Mock
}

func (r *MockProductRepository) GetProduct(ctx context.Context, id uuid.UUID) (*foodji.Product, error) {
	return nil, nil
}

func (r *MockProductRepository) ListProducts(ctx context.Context) ([]uuid.UUID, error) {
	return nil, nil
}

func TestVoteService_CreateOrUpdateVote(t *testing.T) {
	// Arrange
	mockVoteRepo := new(MockVoteRepository)
	mockProductRepo := new(MockProductRepository)
	service := application.NewVoteService(mockVoteRepo, mockProductRepo)
	ctx := context.Background()
	sessionID := uuid.New()
	productID := uuid.New()

	t.Run("create new vote when none exists", func(t *testing.T) {
		// Arrange
		score := 5
		mockVoteRepo.On("GetBySessionID", ctx, sessionID).Return([]*domain.Vote{}, nil).Once()
		mockVoteRepo.On("Create", ctx, mock.AnythingOfType("*domain.Vote")).Return(nil).Once()

		// Act
		vote, err := service.CreateOrUpdateVote(ctx, sessionID, productID, score)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, vote)
		assert.Equal(t, sessionID, vote.SessionID)
		assert.Equal(t, productID, vote.ProductID)
		assert.Equal(t, score, vote.Score)

		mockVoteRepo.AssertExpectations(t)
	})

	t.Run("update existing vote", func(t *testing.T) {
		// Arrange
		existingVote, _ := domain.NewVote(sessionID, productID, 3)
		mockVoteRepo.On("GetBySessionID", ctx, sessionID).Return([]*domain.Vote{existingVote}, nil).Once()
		mockVoteRepo.On("Update", ctx, mock.AnythingOfType("*domain.Vote")).Return(nil).Once()

		newScore := 4

		// Act
		vote, err := service.CreateOrUpdateVote(ctx, sessionID, productID, newScore)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, vote)
		assert.Equal(t, sessionID, vote.SessionID)
		assert.Equal(t, productID, vote.ProductID)
		assert.Equal(t, newScore, vote.Score)

		mockVoteRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		// Arrange
		mockVoteRepo.On("GetBySessionID", ctx, sessionID).Return([]*domain.Vote{}, assert.AnError).Once()

		// Act
		vote, err := service.CreateOrUpdateVote(ctx, sessionID, productID, 5)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, vote)

		mockVoteRepo.AssertExpectations(t)
	})

	t.Run("invalid score", func(t *testing.T) {
		// Act
		vote, err := service.CreateOrUpdateVote(ctx, sessionID, productID, 6)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, domain.ErrInvalidScore, err)
		assert.Nil(t, vote)
	})
}

func TestVoteService_GetVotesBySession(t *testing.T) {
	// Arrange
	mockVoteRepo := new(MockVoteRepository)
	mockProductRepo := new(MockProductRepository)
	service := application.NewVoteService(mockVoteRepo, mockProductRepo)
	ctx := context.Background()
	sessionID := uuid.New()

	t.Run("returns votes from repository", func(t *testing.T) {
		// Arrange
		productID := uuid.New()
		vote1, _ := domain.NewVote(sessionID, productID, 4)
		vote2, _ := domain.NewVote(sessionID, uuid.New(), 5)
		expectedVotes := []*domain.Vote{vote1, vote2}

		mockVoteRepo.On("GetBySessionID", ctx, sessionID).Return(expectedVotes, nil).Once()

		// Act
		votes, err := service.GetVotesBySession(ctx, sessionID)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, expectedVotes, votes)
		assert.Len(t, votes, 2)

		mockVoteRepo.AssertExpectations(t)
	})

	t.Run("handles repository errors", func(t *testing.T) {
		// Arrange
		mockVoteRepo.On("GetBySessionID", ctx, sessionID).Return([]*domain.Vote{}, assert.AnError).Once()

		// Act
		votes, err := service.GetVotesBySession(ctx, sessionID)

		// Assert
		assert.Error(t, err)
		assert.Empty(t, votes)

		mockVoteRepo.AssertExpectations(t)
	})
}

func TestVoteService_GetAggregatedScores(t *testing.T) {
	// Arrange
	mockVoteRepo := new(MockVoteRepository)
	mockProductRepo := new(MockProductRepository)
	service := application.NewVoteService(mockVoteRepo, mockProductRepo)
	ctx := context.Background()

	t.Run("returns aggregated scores from repository", func(t *testing.T) {
		// Arrange
		expectedScores := []*domain.ProductScore{
			{
				ProductID: uuid.New(),
				AvgScore:  4.5,
				VoteCount: 10,
			},
			{
				ProductID: uuid.New(),
				AvgScore:  3.8,
				VoteCount: 5,
			},
		}

		mockVoteRepo.On("GetAggregatedScores", ctx).Return(expectedScores, nil).Once()

		// Act
		scores, err := service.GetAggregatedScores(ctx)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, expectedScores, scores)
		assert.Len(t, scores, 2)

		mockVoteRepo.AssertExpectations(t)
	})

	t.Run("handles repository errors", func(t *testing.T) {
		// Arrange
		mockVoteRepo.On("GetAggregatedScores", ctx).Return([]*domain.ProductScore{}, assert.AnError).Once()

		// Act
		scores, err := service.GetAggregatedScores(ctx)

		// Assert
		assert.Error(t, err)
		assert.Empty(t, scores)

		mockVoteRepo.AssertExpectations(t)
	})
}
