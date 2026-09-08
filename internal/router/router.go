package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/yasersyafa/go-schedule/internal/activity"
)

func New(activityHandler *activity.Handler) *gin.Engine {
	r := gin.Default()

	// CORS setup
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"https://janeismine.netlify.app"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders: []string{"Content-Type"},
		AllowCredentials: false,
	}))

	api := r.Group("/api/v1")
	{
		activities := api.Group("/activities")
		{
			activities.GET("", activityHandler.ListAll)
			activities.POST("", activityHandler.Create)
			activities.PUT("/:id", activityHandler.Update)
			activities.DELETE("/:id", activityHandler.Delete)
		}

		api.GET("/days/:day/activities", activityHandler.ListByDay)
		api.GET("/days/:day/free-slots", activityHandler.ListFreeSlots)
	}

	return r
}