package routes

import (
	"github.com/datarohit/go-restaurant-management-backend-project/controllers"

	"github.com/gin-gonic/gin"
)

func InvoiceRoutes(router *gin.Engine) {
	api := router.Group("/api/v1")
	{
		invoices := api.Group("/invoices")
		{
			invoices.POST("/", controllers.CreateInvoice())
			invoices.GET("/", controllers.GetAllInvoices())
			invoices.GET("/:invoiceId", controllers.GetInvoiceByID())
			invoices.PATCH("/:invoiceId", controllers.UpdateInvoiceByID())
		}
	}
}
