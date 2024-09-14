package routes

import (
	"github.com/datarohit/go-restaurant-management-backend-project/controllers"

	"github.com/gin-gonic/gin"
)

func FoodRoutes(router *gin.Engine) {
	api := router.Group("/api/v1")
	{
		foods := api.Group("/foods")
		{
			foods.POST("/", controllers.CreateFood())
			foods.GET("/", controllers.GetAllFoodItems())
			foods.GET("/:foodId", controllers.GetFoodByID())
		}
	}
}
