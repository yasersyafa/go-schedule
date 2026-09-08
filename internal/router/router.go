package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/yasersyafa/go-schedule/internal/activity"
	"github.com/yasersyafa/go-schedule/internal/auth"
)

func New(activityHandler *activity.Handler, authHandler *auth.Handler, apiToken string) *gin.Engine {
	r := gin.Default()

	// CORS setup
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"https://janeismine.netlify.app"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
		AllowCredentials: false,
	}))

	api := r.Group("/api/v1")
	{
		api.POST("/auth/login", authHandler.Login)
		
		api.GET("/days/:day/activities", activityHandler.ListByDay)
		api.GET("/days/:day/free-slots", activityHandler.ListFreeSlots)
		api.GET("/activities", activityHandler.ListAll)

		protected := api.Group("/activities")
		protected.Use(auth.RequireAuth(apiToken))
		{
			protected.POST("", activityHandler.Create)
			protected.PUT("/:id", activityHandler.Update)
			protected.DELETE("/:id", activityHandler.Delete)
		}
		
	}

	return r
}