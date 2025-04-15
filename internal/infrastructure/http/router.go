package http

import (
	"github.com/ArtemSind/food_tinder/internal/application"
	"github.com/ArtemSind/food_tinder/internal/infrastructure/http/handlers"
	"github.com/gorilla/mux"
)

func NewRouter(
	sessionService *application.SessionService,
	voteService *application.VoteService,
) *mux.Router {
	r := mux.NewRouter()

	// Session handlers
	sessionHandler := handlers.NewSessionHandler(sessionService)
	r.HandleFunc("/api/sessions", sessionHandler.CreateSession).Methods("POST")
	r.HandleFunc("/api/sessions/{sessionID}", sessionHandler.GetSession).Methods("GET")

	// Vote handlers
	voteHandler := handlers.NewVoteHandler(voteService)
	r.HandleFunc("/api/sessions/{sessionID}/votes", voteHandler.CreateOrUpdateVote).Methods("POST")
	r.HandleFunc("/api/sessions/{sessionID}/votes", voteHandler.GetVotesBySession).Methods("GET")
	r.HandleFunc("/api/votes/aggregated", voteHandler.GetAggregatedScores).Methods("GET")

	return r
}
