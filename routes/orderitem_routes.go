package routes

import (
	"github.com/datarohit/go-restaurant-management-backend-project/controllers"

	"github.com/gin-gonic/gin"
)

func OrderItemRoutes(router *gin.Engine) {
	api := router.Group("/api/v1")
	{
		orderItems := api.Group("/orderItems")
		{
			orderItems.POST("/", controllers.CreateOrderItem())
			orderItems.GET("/", controllers.GetAllOrderItems())
			orderItems.GET("/order/:orderId", controllers.GetOrderItemsByOrderID())
			orderItems.GET("/:orderItemId", controllers.GetOrderItemByID())
			orderItems.PATCH("/:orderItemId", controllers.UpdateOrderItemByID())
		}
	}
}
