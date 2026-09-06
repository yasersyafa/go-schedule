package activity

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

var ErrOverlap = errors.New("activity overlaps with an existing one on this day.")

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListAll(ctx context.Context) ([]Activity, error) {
	return s.repo.ListAll(ctx)
}

func (s *Service) ListByDay(ctx context.Context, day string) ([]Activity, error) {
	return s.repo.ListByDay(ctx, day)
}

func (s *Service) Create(ctx context.Context, a *Activity) error {
	overlap, err := s.repo.HasOverlap(ctx, a.Day, a.StartTime, a.EndTime, nil)
	if err != nil {
		return fmt.Errorf("check overlap: %w", err)
	}
	if overlap {
		return ErrOverlap
	}
	return s.repo.Create(ctx, a)
}

func (s *Service) Update(ctx context.Context, a *Activity) error {
	overlap, err := s.repo.HasOverlap(ctx, a.Day, a.StartTime, a.EndTime, &a.ID)
	if err != nil {
		return fmt.Errorf("check overlap: %w", err)
	}
	if overlap {
		return ErrOverlap
	}
	return s.repo.Update(ctx, a)
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}