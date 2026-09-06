package scheduler

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/robfig/cron/v3"
	"github.com/yasersyafa/go-schedule/internal/notifier"
)


type dueActivity struct {
	ID string `db:"id"`
	Name string `db:"name"`
}

type Scheduler struct {
	db *sqlx.DB
	notifier *notifier.TelegramNotifier
	cron *cron.Cron
}

func New(db *sqlx.DB, notifier *notifier.TelegramNotifier) *Scheduler {
	return &Scheduler{
		db: db,
		notifier: notifier,
		cron: cron.New(),
	}
}

func (s *Scheduler) Start() error {
	_, err := s.cron.AddFunc("* * * * *", s.checkDueActivities)
	if err != nil {
		return fmt.Errorf("register cron job: %w", err)
	}
	s.cron.Start()
	log.Println("scheduler started, checking every minute")
	return nil
}

func (s *Scheduler) checkDueActivities() {
	ctx := context.Background()
	now := time.Now()
	currentDay := strings.ToLower(now.Weekday().String())
	currentTime := now.Format("15:04:00")

	var activities []dueActivity
	query := `
		SELECT id, name FROM activities
		WHERE day = $1
		AND start_time::text = $2
		AND (last_notified_date IS NULL OR last_notified_date != CURRENT_DATE)
	`

	if err := s.db.SelectContext(ctx, &activities, query, currentDay, currentTime); err != nil {
		log.Printf("scheduler query error: %v", err)
		return
	}

	for _, a := range activities {
		message := fmt.Sprintf("Hii! Yuk sekarang jadwalnya: %s", a.Name)
		if err := s.notifier.Send(ctx, message); err != nil {
			log.Printf("failed to notify activity %s: %v", a.ID, err)
			continue
		}

		tx, err := s.db.BeginTxx(ctx, nil)
		if err != nil {
			log.Printf("failed to begin tx for activity %s: %v", a.ID, err)
			continue
		}

		if _, err := tx.ExecContext(ctx, `
			UPDATE activities SET last_notified_date = CURRENT_DATE WHERE id = $1
		`, a.ID); err != nil {
			log.Printf("failed to update last_notified_date for %s: %v", a.Name, err)
			tx.Rollback()
			continue
		}

		if _, err := tx.ExecContext(ctx,
			`INSERT INTO notification_logs (activity_id) VALUES ($1)`, a.ID,
		); err != nil {
			log.Printf("failed to insert notification log for %s: %v", a.ID, err)
			tx.Rollback()
			continue
		}

		if err := tx.Commit(); err != nil {
			log.Printf("failed to commit notification for %s: %v", a.ID, err)
		}
	}
}