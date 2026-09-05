package activity

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) ListByDay(ctx context.Context, day string) ([]Activity, error) {
	var activities []Activity
	query := `SELECT * FROM activities WHERE day = $1 ORDER BY start_time`
	if err := r.db.SelectContext(ctx, &activities, query, day); err != nil {
		return nil, fmt.Errorf("list activities by day: %w", err)
	}

	return activities, nil
}

func (r *Repository) HasOverlap(ctx context.Context, day, startTime, endTime string, excludeID *uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(
		SELECT 1 FROM activities WHERE day = $1
		AND start_time < $3
		AND end_time > $2
		AND ($4::uuid IS NULL OR id != $4)
	)`

	var exists bool
	if err := r.db.SelectContext(ctx, &exists, query, day, startTime, endTime, excludeID); err != nil {
		return false, fmt.Errorf("check overlap: %w", err)
	}

	return exists, nil
}

func (r *Repository) Create(ctx context.Context, a *Activity) error {
	query := `INSERT INTO activities (name, notes, day, start_time, end_time)
		VALUES (:name, :notes, :day, :start_time, :end_time)
		RETURNING id, created_at, updated_at
	`

	rows, err := r.db.NamedQueryContext(ctx, query, a)
	if err != nil {
		return fmt.Errorf("create activity: %w", err)
	}

	defer rows.Close()

	if rows.Next() {
		return rows.Scan(&a.ID, &a.CreatedAt, &a.UpdatedAt)
	}

	return nil
}

func (r *Repository) Update(ctx context.Context, a *Activity) error {
	query := `
		UPDATE activities SET name = :name,
		notes = :notes, start_time = :start_time, end_time = :end_time,
		updated_at = now()
		WHERE id = :id
	`

	_, err := r.db.NamedExecContext(ctx, query, a)
	if err != nil {
		return fmt.Errorf("update activity: %w", err)
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM activities WHERE id = :id`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete activity: %w", err)
	}

	return nil
}