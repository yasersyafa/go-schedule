package router

import (
	"github.com/gin-gonic/gin"
	"github.com/yasersyafa/go-schedule/internal/activity"
)

func New(activityHandler *activity.Handler) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")
	{
		activities := api.Group("/activities")
		{
			activities.POST("", activityHandler.Create)
			activities.PUT("/:id", activityHandler.Update)
			activities.DELETE("/:id", activityHandler.Delete)
		}

		api.GET("/days/:day/activities", activityHandler.ListByDay)
	}

	return r
}