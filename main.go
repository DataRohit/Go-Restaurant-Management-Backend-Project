package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/datarohit/go-restaurant-management-backend-project/config"
	"github.com/datarohit/go-restaurant-management-backend-project/middlewares"
	"github.com/datarohit/go-restaurant-management-backend-project/routes"
	"github.com/datarohit/go-restaurant-management-backend-project/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	err := utils.InitializeLogger(zapcore.DebugLevel, []string{"stdout"})
	if err != nil {
		fmt.Printf("Error initializing logger: %v\n", err)
		os.Exit(1)
	}
	log := utils.GetLogger()

	port := config.GetEnvAsInt("PORT", 8080)
	ginMode := config.GetEnv("GIN_MODE", "release")

	gin.SetMode(ginMode)
	router := gin.New()

	router.Use(middlewares.ZapLoggerMiddleware(log))
	router.Use(gin.Recovery())

	routes.HealthRoutes(router)
	routes.UserRoutes(router)

	router.Use(middlewares.Authentication())

	routes.MenuRoutes(router)
	routes.FoodRoutes(router)
	routes.TableRoutes(router)
	routes.OrderRoutes(router)
	routes.OrderItemRoutes(router)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed to start", zap.Error(err))
		}
	}()
	log.Info(fmt.Sprintf("Server is running on port %d", port))

	<-quit
	log.Info("Shutdown signal received, exiting gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown", zap.Error(err))
	}

	log.Info("Server exited cleanly")
}
