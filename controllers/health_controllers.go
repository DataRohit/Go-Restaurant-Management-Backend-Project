package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/datarohit/go-restaurant-management-backend-project/database"
	"github.com/gin-gonic/gin"
)

func GetRouterHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "Router is healthy",
	})
}

func GetDatabaseHealth(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbClient := database.Client
	if dbClient == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "Database client uninitialized",
		})
		return
	}

	if err := dbClient.Ping(ctx, nil); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "Database unreachable",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "Database is healthy",
	})
}
