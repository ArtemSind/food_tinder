package handlers_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ArtemSind/food_tinder/internal/domain"
	"github.com/ArtemSind/food_tinder/internal/infrastructure/http/handlers"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Mock SessionService
type MockSessionService struct {
	mock.Mock
}

func (m *MockSessionService) CreateSession(ctx context.Context) (*domain.Session, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Session), args.Error(1)
}

func (m *MockSessionService) GetSession(ctx context.Context, id uuid.UUID) (*domain.Session, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Session), args.Error(1)
}

func TestSessionHandler_CreateSession(t *testing.T) {
	// Arrange
	mockService := new(MockSessionService)
	handler := handlers.NewSessionHandler(mockService)

	t.Run("successful session creation", func(t *testing.T) {
		// Arrange
		sessionID := uuid.New()
		session := &domain.Session{
			ID: sessionID,
		}
		mockService.On("CreateSession", mock.Anything).Return(session, nil).Once()

		// Create request and recorder
		req := httptest.NewRequest("POST", "/api/sessions", nil)
		rec := httptest.NewRecorder()

		// Act
		handler.CreateSession(rec, req)

		// Assert
		require.Equal(t, http.StatusCreated, rec.Code)

		var responseSession domain.Session
		err := json.NewDecoder(rec.Body).Decode(&responseSession)
		require.NoError(t, err)
		assert.Equal(t, sessionID, responseSession.ID)

		mockService.AssertExpectations(t)
	})

	t.Run("service error", func(t *testing.T) {
		// Arrange
		mockService.On("CreateSession", mock.Anything).Return(nil, errors.New("service error")).Once()

		// Create request and recorder
		req := httptest.NewRequest("POST", "/api/sessions", nil)
		rec := httptest.NewRecorder()

		// Act
		handler.CreateSession(rec, req)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockService.AssertExpectations(t)
	})
}

func TestSessionHandler_GetSession(t *testing.T) {
	// Arrange
	mockService := new(MockSessionService)
	handler := handlers.NewSessionHandler(mockService)

	t.Run("successful session retrieval", func(t *testing.T) {
		// Arrange
		sessionID := uuid.New()
		session := &domain.Session{
			ID: sessionID,
		}
		mockService.On("GetSession", mock.Anything, sessionID).Return(session, nil).Once()

		// Create request and recorder
		req := httptest.NewRequest("GET", "/api/sessions/"+sessionID.String(), nil)
		rec := httptest.NewRecorder()

		// Setup router to get URL params
		router := mux.NewRouter()
		router.HandleFunc("/api/sessions/{sessionID}", handler.GetSession)
		router.ServeHTTP(rec, req)

		// Assert
		require.Equal(t, http.StatusOK, rec.Code)

		var responseSession domain.Session
		err := json.NewDecoder(rec.Body).Decode(&responseSession)
		require.NoError(t, err)
		assert.Equal(t, sessionID, responseSession.ID)

		mockService.AssertExpectations(t)
	})

	t.Run("invalid session ID", func(t *testing.T) {
		// Create request and recorder with invalid UUID
		req := httptest.NewRequest("GET", "/api/sessions/invalid-uuid", nil)
		rec := httptest.NewRecorder()

		// Setup router to get URL params
		router := mux.NewRouter()
		router.HandleFunc("/api/sessions/{sessionID}", handler.GetSession)
		router.ServeHTTP(rec, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("session not found", func(t *testing.T) {
		// Arrange
		sessionID := uuid.New()
		mockService.On("GetSession", mock.Anything, sessionID).Return(nil, errors.New("not found")).Once()

		// Create request and recorder
		req := httptest.NewRequest("GET", "/api/sessions/"+sessionID.String(), nil)
		rec := httptest.NewRecorder()

		// Setup router to get URL params
		router := mux.NewRouter()
		router.HandleFunc("/api/sessions/{sessionID}", handler.GetSession)
		router.ServeHTTP(rec, req)

		// Assert
		assert.Equal(t, http.StatusNotFound, rec.Code)
		mockService.AssertExpectations(t)
	})
}
