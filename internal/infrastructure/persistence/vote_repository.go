package persistence

import (
	"context"
	"fmt"

	"github.com/ArtemSind/food_tinder/internal/domain"
	"github.com/ArtemSind/food_tinder/internal/infrastructure/persistence/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type VoteRepository struct {
	db *pgxpool.Pool
}

func NewVoteRepository(db *pgxpool.Pool) *VoteRepository {
	return &VoteRepository{db: db}
}

func (r *VoteRepository) Create(ctx context.Context, vote *domain.Vote) error {
	dbVote := models.VoteFromDomain(vote)
	_, err := r.db.Exec(ctx,
		"INSERT INTO votes (id, session_id, product_id, score, created_at) VALUES ($1, $2, $3, $4, $5)",
		dbVote.ID, dbVote.SessionID, dbVote.ProductID, dbVote.Score, dbVote.CreatedAt)
	return err
}

func (r *VoteRepository) Update(ctx context.Context, vote *domain.Vote) error {
	dbVote := models.VoteFromDomain(vote)
	_, err := r.db.Exec(ctx,
		"UPDATE votes SET score = $1 WHERE session_id = $2 AND product_id = $3",
		dbVote.Score, dbVote.SessionID, dbVote.ProductID)
	return err
}

func (r *VoteRepository) GetBySessionID(ctx context.Context, sessionID uuid.UUID) ([]*domain.Vote, error) {
	rows, err := r.db.Query(ctx,
		`SELECT v.id, v.session_id, v.product_id, v.score, v.created_at
		FROM votes v
		WHERE v.session_id = $1`, sessionID)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	var dbVotes []*models.VoteDB
	for rows.Next() {
		var dbVote models.VoteDB
		if err := rows.Scan(&dbVote.ID, &dbVote.SessionID, &dbVote.ProductID, &dbVote.Score, &dbVote.CreatedAt); err != nil {
			return nil, err
		}
		dbVotes = append(dbVotes, &dbVote)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return models.VotesToDomain(dbVotes), nil
}

func (r *VoteRepository) GetAggregatedScores(ctx context.Context) ([]*domain.ProductScore, error) {
	rows, err := r.db.Query(ctx,
		`SELECT 
			product_id, 
			AVG(score) as avg_score, 
			COUNT(id) as vote_count
		FROM votes
		GROUP BY product_id
		ORDER BY avg_score DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scores []*domain.ProductScore
	for rows.Next() {
		score := &domain.ProductScore{}
		if err := rows.Scan(&score.ProductID, &score.AvgScore, &score.VoteCount); err != nil {
			return nil, err
		}
		scores = append(scores, score)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return scores, nil
}
