package application

import (
	"context"

	"github.com/ArtemSind/food_tinder/internal/domain"
	"github.com/google/uuid"
)

type SessionRepository interface {
	Create(ctx context.Context, session *domain.Session) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Session, error)
}

type SessionService struct {
	repo SessionRepository
}

func NewSessionService(repo SessionRepository) *SessionService {
	return &SessionService{repo: repo}
}

func (s *SessionService) CreateSession(ctx context.Context) (*domain.Session, error) {
	session := domain.NewSession()
	if err := s.repo.Create(ctx, session); err != nil {
		return nil, err
	}
	return session, nil
}

func (s *SessionService) GetSession(ctx context.Context, id uuid.UUID) (*domain.Session, error) {
	return s.repo.GetByID(ctx, id)
}
