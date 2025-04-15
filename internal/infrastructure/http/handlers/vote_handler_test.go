package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ArtemSind/food_tinder/internal/domain"
	"github.com/ArtemSind/food_tinder/internal/infrastructure/http/handlers"
	"github.com/ArtemSind/food_tinder/internal/infrastructure/http/models"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Mock VoteService
type MockVoteService struct {
	mock.Mock
}

func (m *MockVoteService) CreateOrUpdateVote(ctx context.Context, sessionID, productID uuid.UUID, score int) (*domain.Vote, error) {
	args := m.Called(ctx, sessionID, productID, score)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Vote), args.Error(1)
}

func (m *MockVoteService) GetVotesBySession(ctx context.Context, sessionID uuid.UUID) ([]*domain.Vote, error) {
	args := m.Called(ctx, sessionID)
	return args.Get(0).([]*domain.Vote), args.Error(1)
}

func (m *MockVoteService) GetAggregatedScores(ctx context.Context) ([]*domain.ProductScore, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*domain.ProductScore), args.Error(1)
}

func TestVoteHandler_CreateOrUpdateVote(t *testing.T) {
	// Arrange
	mockService := new(MockVoteService)
	handler := handlers.NewVoteHandler(mockService)

	t.Run("successful vote creation", func(t *testing.T) {
		// Arrange
		sessionID := uuid.New()
		productID := uuid.New()
		score := 4

		vote := &domain.Vote{
			ID:        uuid.New(),
			SessionID: sessionID,
			ProductID: productID,
			Score:     score,
		}

		requestBody, _ := json.Marshal(map[string]interface{}{
			"product_id": productID.String(),
			"score":      score,
		})

		mockService.On("CreateOrUpdateVote", mock.Anything, sessionID, productID, score).Return(vote, nil).Once()

		// Create request and recorder
		req := httptest.NewRequest("POST", "/api/sessions/"+sessionID.String()+"/votes", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		// Setup router to get URL params
		router := mux.NewRouter()
		router.HandleFunc("/api/sessions/{sessionID}/votes", handler.CreateOrUpdateVote).Methods("POST")
		router.ServeHTTP(rec, req)

		// Assert
		require.Equal(t, http.StatusOK, rec.Code)

		var responseVote domain.Vote
		err := json.NewDecoder(rec.Body).Decode(&responseVote)
		require.NoError(t, err)
		assert.Equal(t, vote.ID, responseVote.ID)
		assert.Equal(t, sessionID, responseVote.SessionID)
		assert.Equal(t, productID, responseVote.ProductID)
		assert.Equal(t, score, responseVote.Score)

		mockService.AssertExpectations(t)
	})

	t.Run("invalid session ID", func(t *testing.T) {
		// Create request with invalid session ID
		requestBody, _ := json.Marshal(map[string]interface{}{
			"product_id": uuid.New().String(),
			"score":      4,
		})

		req := httptest.NewRequest("POST", "/api/sessions/invalid-uuid/votes", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		// Setup router to get URL params
		router := mux.NewRouter()
		router.HandleFunc("/api/sessions/{sessionID}/votes", handler.CreateOrUpdateVote).Methods("POST")
		router.ServeHTTP(rec, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("invalid product ID", func(t *testing.T) {
		// Arrange
		sessionID := uuid.New()

		requestBody, _ := json.Marshal(map[string]interface{}{
			"product_id": "invalid-uuid",
			"score":      4,
		})

		req := httptest.NewRequest("POST", "/api/sessions/"+sessionID.String()+"/votes", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		// Setup router to get URL params
		router := mux.NewRouter()
		router.HandleFunc("/api/sessions/{sessionID}/votes", handler.CreateOrUpdateVote).Methods("POST")
		router.ServeHTTP(rec, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("invalid JSON body", func(t *testing.T) {
		// Arrange
		sessionID := uuid.New()

		// Invalid JSON
		requestBody := []byte(`{"product_id": "invalid-json`)

		req := httptest.NewRequest("POST", "/api/sessions/"+sessionID.String()+"/votes", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		// Setup router to get URL params
		router := mux.NewRouter()
		router.HandleFunc("/api/sessions/{sessionID}/votes", handler.CreateOrUpdateVote).Methods("POST")
		router.ServeHTTP(rec, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("service error", func(t *testing.T) {
		// Arrange
		sessionID := uuid.New()
		productID := uuid.New()
		score := 4

		requestBody, _ := json.Marshal(map[string]interface{}{
			"product_id": productID.String(),
			"score":      score,
		})

		mockService.On("CreateOrUpdateVote", mock.Anything, sessionID, productID, score).
			Return(nil, errors.New("service error")).Once()

		// Create request and recorder
		req := httptest.NewRequest("POST", "/api/sessions/"+sessionID.String()+"/votes", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		// Setup router to get URL params
		router := mux.NewRouter()
		router.HandleFunc("/api/sessions/{sessionID}/votes", handler.CreateOrUpdateVote).Methods("POST")
		router.ServeHTTP(rec, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		mockService.AssertExpectations(t)
	})
}

func TestVoteHandler_GetVotesBySession(t *testing.T) {
	// Arrange
	mockService := new(MockVoteService)
	handler := handlers.NewVoteHandler(mockService)

	t.Run("successful votes retrieval", func(t *testing.T) {
		// Arrange
		sessionID := uuid.New()
		productID1 := uuid.New()
		productID2 := uuid.New()

		vote1, _ := domain.NewVote(sessionID, productID1, 4)
		vote2, _ := domain.NewVote(sessionID, productID2, 5)
		expectedVotes := []*domain.Vote{vote1, vote2}

		mockService.On("GetVotesBySession", mock.Anything, sessionID).Return(expectedVotes, nil).Once()

		// Create request and recorder
		req := httptest.NewRequest("GET", "/api/sessions/"+sessionID.String()+"/votes", nil)
		rec := httptest.NewRecorder()

		// Setup router to get URL params
		router := mux.NewRouter()
		router.HandleFunc("/api/sessions/{sessionID}/votes", handler.GetVotesBySession).Methods("GET")
		router.ServeHTTP(rec, req)

		// Assert
		require.Equal(t, http.StatusOK, rec.Code)

		var responseData struct {
			Votes []*models.VoteResponse `json:"votes"`
			Count int                    `json:"count"`
		}
		err := json.NewDecoder(rec.Body).Decode(&responseData)
		require.NoError(t, err)
		assert.Len(t, responseData.Votes, 2)
		assert.Equal(t, expectedVotes[0].ID, responseData.Votes[0].ID)
		assert.Equal(t, expectedVotes[1].ID, responseData.Votes[1].ID)

		mockService.AssertExpectations(t)
	})

	t.Run("invalid session ID", func(t *testing.T) {
		// Create request with invalid session ID
		req := httptest.NewRequest("GET", "/api/sessions/invalid-uuid/votes", nil)
		rec := httptest.NewRecorder()

		// Setup router to get URL params
		router := mux.NewRouter()
		router.HandleFunc("/api/sessions/{sessionID}/votes", handler.GetVotesBySession).Methods("GET")
		router.ServeHTTP(rec, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("service error", func(t *testing.T) {
		// Arrange
		sessionID := uuid.New()
		mockService.On("GetVotesBySession", mock.Anything, sessionID).
			Return([]*domain.Vote{}, errors.New("service error")).Once()

		// Create request and recorder
		req := httptest.NewRequest("GET", "/api/sessions/"+sessionID.String()+"/votes", nil)
		rec := httptest.NewRecorder()

		// Setup router to get URL params
		router := mux.NewRouter()
		router.HandleFunc("/api/sessions/{sessionID}/votes", handler.GetVotesBySession).Methods("GET")
		router.ServeHTTP(rec, req)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockService.AssertExpectations(t)
	})
}

func TestVoteHandler_GetAggregatedScores(t *testing.T) {
	// Arrange
	mockService := new(MockVoteService)
	handler := handlers.NewVoteHandler(mockService)

	t.Run("successful aggregated scores retrieval", func(t *testing.T) {
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

		mockService.On("GetAggregatedScores", mock.Anything).Return(expectedScores, nil).Once()

		// Create request and recorder
		req := httptest.NewRequest("GET", "/api/votes/aggregated", nil)
		rec := httptest.NewRecorder()

		// Act
		handler.GetAggregatedScores(rec, req)

		// Assert
		require.Equal(t, http.StatusOK, rec.Code)

		var responseData struct {
			Scores []*models.ProductScoreResponse `json:"scores"`
			Count  int                            `json:"count"`
		}
		err := json.NewDecoder(rec.Body).Decode(&responseData)
		require.NoError(t, err)
		assert.Len(t, responseData.Scores, 2)
		assert.Equal(t, expectedScores[0].ProductID, responseData.Scores[0].ProductID)
		assert.Equal(t, expectedScores[0].AvgScore, responseData.Scores[0].AvgScore)
		assert.Equal(t, expectedScores[0].VoteCount, responseData.Scores[0].VoteCount)

		mockService.AssertExpectations(t)
	})

	t.Run("service error", func(t *testing.T) {
		// Arrange
		mockService.On("GetAggregatedScores", mock.Anything).
			Return([]*domain.ProductScore{}, errors.New("service error")).Once()

		// Create request and recorder
		req := httptest.NewRequest("GET", "/api/votes/aggregated", nil)
		rec := httptest.NewRecorder()

		// Act
		handler.GetAggregatedScores(rec, req)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockService.AssertExpectations(t)
	})
}
