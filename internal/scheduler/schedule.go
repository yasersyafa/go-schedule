package scheduler

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
	"github.com/yasersyafa/go-schedule/internal/activity"
	"github.com/yasersyafa/go-schedule/internal/notifier"
)

type schedulerRepository interface {
	ListDueNow(ctx context.Context, day, currentTime, currentDate string)([]activity.DueActivity, error)
	RecordNotification(ctx context.Context, activityID uuid.UUID, notifiedDate string) error
}

type Scheduler struct {
	repo schedulerRepository
	notifiers []notifier.Notifier
	cron *cron.Cron
	loc *time.Location
}

func New(repo schedulerRepository, tz string, notifiers ...notifier.Notifier) (*Scheduler, error) {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return nil, fmt.Errorf("load timezone %q: %w", tz, err)
	}
	return &Scheduler{
		repo: repo,
		notifiers: notifiers,
		cron: cron.New(cron.WithLocation(loc)),
		loc: loc,
	}, nil
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
	now := time.Now().In(s.loc)
	currentDay := strings.ToLower(now.Weekday().String())
	currentTime := now.Format("15:04:05")
	currentDate := now.Format("2006-01-02")

	activities, err := s.repo.ListDueNow(ctx, currentDay, currentTime, currentDate)
	if err != nil {
		log.Printf("scheduler query error: %v", err)
		return
	}

	for _, a := range activities {
		message := fmt.Sprintf("sekarang jadwalnya %s nih! semangat yaa", a.Name)

		sentToAny := false

		for _, n := range s.notifiers {
			if err := n.Send(ctx, message); err != nil {
				log.Printf("failed to notify activity %s via one channel: %v", a.ID, err)
				continue
			}
			sentToAny = true
		}

		if !sentToAny {
			continue
		}

		if err := s.repo.RecordNotification(ctx, a.ID, currentDate); err != nil {
			log.Printf("failed to to record notification for %s: %v", a.ID, err)
		}
	}
}