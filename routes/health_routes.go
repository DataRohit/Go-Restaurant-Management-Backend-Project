package routes

import (
	"github.com/datarohit/go-restaurant-management-backend-project/controllers"
	"github.com/gin-gonic/gin"
)

func HealthRoutes(router *gin.Engine) {
	health := router.Group("/health")
	{
		health.GET("/router", controllers.GetRouterHealth)
		health.GET("/database", controllers.GetDatabaseHealth)
	}
}
