package activity

import (
	"time"

	"github.com/google/uuid"
)

type Activity struct {
	ID uuid.UUID `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
	Notes *string `db:"notes" json:"notes"`
	Day string `db:"day" json:"day"`
	StartTime string `db:"start_time" json:"start_time"`
	EndTime string `db:"end_time" json:"end_time"`
	LastNotifiedDate *time.Time `db:"last_notified_date" json:"last_notified_date"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type FreeSlot struct {
	Start string `json:"start"`
	End string `json:"end"`
}

type DueActivity struct {
	ID uuid.UUID `json:"id"`
	Name string `json:"name"`
}