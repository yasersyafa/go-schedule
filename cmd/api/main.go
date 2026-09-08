package main

import (
	"log"

	_ "time/tzdata"

	"github.com/yasersyafa/go-schedule/internal/activity"
	"github.com/yasersyafa/go-schedule/internal/auth"
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

	// auth
	authHandler := auth.NewHandler(cfg.AdminPassword, cfg.ApiToken)

	tgNotifier := notifier.NewTelegramNotifier(cfg.TelegramBotToken, cfg.TelegramChatID)
	sched, err := scheduler.New(conn, tgNotifier, cfg.Timezone)
	if err != nil {
		log.Fatalf("failed to init scheduler: %v", err)
	}
	if err := sched.Start(); err != nil {
		log.Fatalf("failed to start scheduler: %v", err)
	}

	r := router.New(activityHandler, authHandler, cfg.ApiToken)

	log.Printf("server starting on port: %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}