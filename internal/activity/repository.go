package activity

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var activityColumns string = `
	id, name, notes, day, start_time::text AS start_time,
	end_time::text AS end_time, last_notified_date, created_at, updated_at
`

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) ListAll(ctx context.Context) ([]Activity, error) {
	activities := []Activity{}
	query := `
		SELECT ` + activityColumns  + ` FROM activities
		ORDER BY
			CASE day
				WHEN 'monday' THEN 1
				WHEN 'tuesday' THEN 2
				WHEN 'wednesday' THEN 3
				WHEN 'thursday' THEN 4
				WHEN 'friday' THEN 5
				WHEN 'saturday' THEN 6
				WHEN 'sunday' THEN 7
			END,
			start_time
	`

	if err := r.db.SelectContext(ctx, &activities, query); err != nil {
		return nil, fmt.Errorf("list all activities: %w", err)
	}
	return activities, nil
}

func (r *Repository) ListByDay(ctx context.Context, day string) ([]Activity, error) {
	activities := []Activity{}
	query := `SELECT ` + activityColumns + ` FROM activities WHERE day = $1 ORDER BY start_time`
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
	if err := r.db.GetContext(ctx, &exists, query, day, startTime, endTime, excludeID); err != nil {
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
	query := `DELETE FROM activities WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete activity: %w", err)
	}

	return nil
}

func (r *Repository) ListDueNow(ctx context.Context, day, currentTime, currentDate string) ([]DueActivity, error) {
	var activities []DueActivity

	query := `
		SELECT id, name FROM activities
		WHERE day = $1
		AND start_time::text = $2
		AND (last_notified_date IS NULL OR last_notified_date != $3)
	`

	if err := r.db.SelectContext(ctx, &activities, query, day, currentTime, currentDate); err != nil {
		return []DueActivity{}, fmt.Errorf("list due activity: %w", err)
	}

	return activities, nil
}

func (r *Repository) RecordNotification(ctx context.Context, activityID uuid.UUID, notifiedDate string) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	if _, err := tx.ExecContext(ctx,
		`UPDATE activities SET last_notified_date = $2 WHERE id = $1`, activityID, notifiedDate,
	); err != nil {
		tx.Rollback()
		return fmt.Errorf("update last_notified_date: %w", err)
	}

	if _, err := tx.ExecContext(ctx,
		`INSERT INTO notif_logs (activity_id) VALUES ($1)`, activityID,
	); err != nil {
		tx.Rollback()
		return fmt.Errorf("insert notification log: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	return nil
}