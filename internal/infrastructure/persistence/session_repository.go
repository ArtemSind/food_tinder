package persistence

import (
	"context"

	"github.com/ArtemSind/food_tinder/internal/domain"
	"github.com/ArtemSind/food_tinder/internal/infrastructure/persistence/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type SessionRepository struct {
	db *pgxpool.Pool
}

func NewSessionRepository(db *pgxpool.Pool) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) Create(ctx context.Context, session *domain.Session) error {
	dbSession := models.SessionFromDomain(session)
	_, err := r.db.Exec(ctx,
		"INSERT INTO sessions (id, created_at) VALUES ($1, $2)",
		dbSession.ID, dbSession.CreatedAt)
	return err
}

func (r *SessionRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Session, error) {
	var dbSession models.SessionDB
	err := r.db.QueryRow(ctx,
		"SELECT id, created_at FROM sessions WHERE id = $1", id).
		Scan(&dbSession.ID, &dbSession.CreatedAt)
	if err != nil {
		return nil, err
	}
	return dbSession.ToDomain(), nil
}
