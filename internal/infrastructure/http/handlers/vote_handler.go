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

type VoteService interface {
	CreateOrUpdateVote(ctx context.Context, sessionID, productID uuid.UUID, score int) (*domain.Vote, error)
	GetVotesBySession(ctx context.Context, sessionID uuid.UUID) ([]*domain.Vote, error)
	GetAggregatedScores(ctx context.Context) ([]*domain.ProductScore, error)
}

type VoteHandler struct {
	voteService VoteService
}

func NewVoteHandler(voteService VoteService) *VoteHandler {
	return &VoteHandler{voteService: voteService}
}

func (h *VoteHandler) CreateOrUpdateVote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID, err := uuid.Parse(vars["sessionID"])
	if err != nil {
		http.Error(w, "Invalid session ID", http.StatusBadRequest)
		return
	}

	var req httpModels.VoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	vote, err := req.ToDomain(sessionID)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	vote, err = h.voteService.CreateOrUpdateVote(ctx, vote.SessionID, vote.ProductID, vote.Score)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := httpModels.VoteResponseFromDomain(vote)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *VoteHandler) GetVotesBySession(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID, err := uuid.Parse(vars["sessionID"])
	if err != nil {
		http.Error(w, "Invalid session ID", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	votes, err := h.voteService.GetVotesBySession(ctx, sessionID)
	if err != nil {
		http.Error(w, "Failed to get votes", http.StatusInternalServerError)
		return
	}

	response := httpModels.VoteListResponseFromDomain(votes)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *VoteHandler) GetAggregatedScores(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	scores, err := h.voteService.GetAggregatedScores(ctx)
	if err != nil {
		http.Error(w, "Failed to get aggregated scores", http.StatusInternalServerError)
		return
	}

	response := httpModels.ProductScoreListResponseFromDomain(scores)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
