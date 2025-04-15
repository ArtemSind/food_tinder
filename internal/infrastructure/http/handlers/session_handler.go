package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ArtemSind/food_tinder/internal/domain"
	httpModels "github.com/ArtemSind/food_tinder/internal/infrastructure/http/models"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type SessionService interface {
	CreateSession(ctx context.Context) (*domain.Session, error)
	GetSession(ctx context.Context, id uuid.UUID) (*domain.Session, error)
}

type SessionHandler struct {
	sessionService SessionService
}

func NewSessionHandler(sessionService SessionService) *SessionHandler {
	return &SessionHandler{sessionService: sessionService}
}

func (h *SessionHandler) CreateSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	session, err := h.sessionService.CreateSession(ctx)
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	response := httpModels.SessionResponseFromDomain(session)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *SessionHandler) GetSession(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID, err := uuid.Parse(vars["sessionID"])
	if err != nil {
		http.Error(w, "Invalid session ID", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	session, err := h.sessionService.GetSession(ctx, sessionID)
	if err != nil {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	response := httpModels.SessionResponseFromDomain(session)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
