package routes

import (
	"github.com/datarohit/go-restaurant-management-backend-project/controllers"

	"github.com/gin-gonic/gin"
)

func TableRoutes(router *gin.Engine) {
	api := router.Group("/api/v1")
	{
		tables := api.Group("/tables")
		{
			tables.POST("/", controllers.CreateTable())
			tables.GET("/", controllers.GetAllTables())
			tables.GET("/:tableId", controllers.GetTableByID())
			tables.PATCH("/:tableId", controllers.UpdateTableByID())
		}
	}
}
