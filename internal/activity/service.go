package activity

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

const (
	dayStart = "00:00:00"
	dayEnd = "23:59:59"
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

func (s *Service) ListFreeSlots(ctx context.Context, day string) ([]FreeSlot, error) {
	activities, err := s.repo.ListByDay(ctx, day)
	if err != nil {
		return nil, fmt.Errorf("list activities for free slots: %w", err)
	}

	var freeSlots []FreeSlot
	cursor := dayStart

	for _, a := range activities {
		if a.StartTime > cursor {
			freeSlots = append(freeSlots, FreeSlot{ Start: cursor, End: a.StartTime })
		}
		if a.EndTime > cursor {
			cursor = a.EndTime
		}
	}

	if cursor < dayEnd {
		freeSlots = append(freeSlots, FreeSlot{ Start: cursor, End: dayEnd })
	}

	return freeSlots, nil
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