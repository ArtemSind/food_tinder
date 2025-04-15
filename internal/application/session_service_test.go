package application_test

import (
	"context"
	"testing"

	"github.com/ArtemSind/food_tinder/internal/application"
	"github.com/ArtemSind/food_tinder/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Mock SessionRepository
type MockSessionRepository struct {
	mock.Mock
}

func (m *MockSessionRepository) Create(ctx context.Context, session *domain.Session) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}

func (m *MockSessionRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Session, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Session), args.Error(1)
}

func TestSessionService_CreateSession(t *testing.T) {
	// Arrange
	mockRepo := new(MockSessionRepository)
	service := application.NewSessionService(mockRepo)
	ctx := context.Background()

	t.Run("creates session successfully", func(t *testing.T) {
		// Arrange
		mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.Session")).Return(nil).Once()

		// Act
		session, err := service.CreateSession(ctx)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, session)
		assert.NotEqual(t, uuid.Nil, session.ID)
		assert.False(t, session.CreatedAt.IsZero())

		mockRepo.AssertExpectations(t)
	})

	t.Run("handles repository error", func(t *testing.T) {
		// Arrange
		mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.Session")).Return(assert.AnError).Once()

		// Act
		session, err := service.CreateSession(ctx)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, session)

		mockRepo.AssertExpectations(t)
	})
}

func TestSessionService_GetSession(t *testing.T) {
	// Arrange
	mockRepo := new(MockSessionRepository)
	service := application.NewSessionService(mockRepo)
	ctx := context.Background()
	sessionID := uuid.New()

	t.Run("gets session by ID", func(t *testing.T) {
		// Arrange
		expectedSession := &domain.Session{
			ID: sessionID,
		}
		mockRepo.On("GetByID", ctx, sessionID).Return(expectedSession, nil).Once()

		// Act
		session, err := service.GetSession(ctx, sessionID)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, expectedSession, session)

		mockRepo.AssertExpectations(t)
	})

	t.Run("handles session not found", func(t *testing.T) {
		// Arrange
		mockRepo.On("GetByID", ctx, sessionID).Return(nil, assert.AnError).Once()

		// Act
		session, err := service.GetSession(ctx, sessionID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, session)

		mockRepo.AssertExpectations(t)
	})
}
