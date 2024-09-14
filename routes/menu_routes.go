package routes

import (
	"github.com/datarohit/go-restaurant-management-backend-project/controllers"

	"github.com/gin-gonic/gin"
)

func MenuRoutes(router *gin.Engine) {
	api := router.Group("/api/v1")
	{
		menus := api.Group("/menus")
		{
			menus.POST("/", controllers.CreateMenu())
			menus.GET("/", controllers.GetAllMenus())
		}
	}
}
