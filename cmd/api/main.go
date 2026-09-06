package main

import (
	"log"

	"github.com/yasersyafa/go-schedule/internal/activity"
	"github.com/yasersyafa/go-schedule/internal/config"
	"github.com/yasersyafa/go-schedule/internal/database"
	"github.com/yasersyafa/go-schedule/internal/notifier"
	"github.com/yasersyafa/go-schedule/internal/router"
	"github.com/yasersyafa/go-schedule/internal/scheduler"
)

func main() {
	cfg := config.Load()

	conn, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	defer conn.Close()

	activityRepo := activity.NewRepository(conn)
	activityService := activity.NewService(activityRepo)
	activityHandler := activity.NewHandler(activityService)

	tgNotifier := notifier.NewTelegramNotifier(cfg.TelegramBotToken, cfg.TelegramChatID)
	sched := scheduler.New(conn, tgNotifier)
	if err := sched.Start(); err != nil {
		log.Fatalf("failed to start scheduler: %v", err)
	}

	r := router.New(activityHandler)

	log.Printf("server starting on port: %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}