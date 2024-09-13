package routes

import (
	"github.com/datarohit/go-restaurant-management-backend-project/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {
	api := router.Group("/api/v1")
	{
		users := api.Group("/users")
		{
			users.POST("/signup", controllers.SignUp())
			users.POST("/login", controllers.Login())
		}
	}
}
