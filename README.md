# Reminder System: GO Schedule

Backend API for a personal weekly study schedule, with real-time Telegram notifications when an activity starts.

## Tech Stack

- Go
- Gin
- PostgreSQL
- sqlx
- golang-migrate
- robfig/cron
- Telegram Bot API
- godotenv
- go-playground/validator.

## Notifications

A cron job runs every minute, checking for activities whose `start_time` matches the current time and haven't been notified today (`last_notified_date`). On a match:

1. Send a message via Telegram Bot API
2. On success, update `last_notified_date` and insert into `notification_logs` within a single transaction, keeping both in sync

Accuracy is within ~59 seconds of the exact time, since it relies on per-minute polling rather than individual timers per activity.

## Running Locally

```bash
# 1. start postgres via docker

docker compose up -d

# 2. Run migrations

migrate -path migrations -database "$DATABASE_URL" up

# 3. Run the server

go run cmd/api/main.go
```

## API Endpoints

| Method     | Path                          | Description                                |
| ---------- | ----------------------------- | ------------------------------------------ |
| \`POST\`   | \`/api/activities\`           | Create a new activity                      |
| \`PUT\`    | \`/api/activities/:id\`       | Edit an activity                           |
| \`DELETE\` | \`/api/activities/:id\`       | Delete an activity                         |
| \`GET\`    | \`/api/activities\`           | List all activities, ordered Monday→Sunday |
| \`GET\`    | \`/api/days/:day/activities\` | List activities for a given day            |
| \`GET\`    | \`/api/days/:day/free-slots\` | List free time slots for a given day       |
