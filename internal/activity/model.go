package activity

import (
	"time"

	"github.com/google/uuid"
)

type Activity struct {
	ID uuid.UUID `db:"id"`
	Name string `db:"name"`
	Notes *string `db:"notes"`
	Day string `db:"day"`
	StartTime string `db:"start_time"`
	EndTime string `db:"end_time"`
	LastNotifiedDate *time.Time `db:"last_notified_date"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type FreeSlot struct {
	Start string `json:"start"`
	End string `json:"end"`
}