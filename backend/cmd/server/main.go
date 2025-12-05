package main

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"

	orderusecase "foodie/backend/internal/application/usecase/order"
	productusecase "foodie/backend/internal/application/usecase/product"
	"foodie/backend/internal/domain/order"
	"foodie/backend/internal/domain/product"
	"foodie/backend/internal/interfaces/http/controller"
	"foodie/backend/internal/interfaces/http/router"
	"foodie/backend/pkg/config"
	"foodie/backend/pkg/logger"

	"go.uber.org/zap"
)

func main() {
	// Load environment variables
	if err := config.Load(); err != nil {
		panic("failed to load config: " + err.Error())
	}

	// Initialize structured logger with JSON output (Grafana/Loki compatible)
	logConfig := logger.Config{
		Level:  config.Get("LOG_LEVEL", "info"),
		Format: config.Get("LOG_FORMAT", "json"),
		Output: config.Get("LOG_OUTPUT", "stdout"),
	}
	appLogger, err := logger.New(logConfig)
	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}
	defer appLogger.Sync()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// TODO: Initialize repositories from infrastructure layer
	// For now, create nil repositories - this will be wired properly later
	// repos, err := database.NewRepositories(sqlDB)
	// if err != nil { ... }

	// TODO: Initialize repositories from infrastructure layer
	// db, err := database.NewConnectionFromEnv()
	// if err != nil { ... }
	// repos, err := database.NewRepositories(db)
	// if err != nil { ... }

	// TODO: Initialize use cases with repositories
	// orderUseCase := orderusecase.NewUseCase(repos.Order, repos.Product)
	// productUseCase := productusecase.NewUseCase(repos.Product)

	// For now, create use cases with nil repositories (will fail at runtime if called)
	// This is a placeholder until full DI is implemented
	var orderRepo order.Repository = nil
	var productRepo product.Repository = nil
	orderUseCase := orderusecase.NewUseCase(orderRepo, productRepo)
	productUseCase := productusecase.NewUseCase(productRepo)

	// Initialize controllers
	healthController := controller.NewHealthController()
	orderController := controller.NewOrderController(orderUseCase)
	productController := controller.NewProductController(productUseCase)

	// Setup router with logger and controllers
	httpRouter := router.NewRouter(appLogger, healthController, orderController, productController)
	httpRouter.SetupRoutes()

	// Server address - can use SERVER_ADDR or combine SERVER_HOST + SERVER_PORT
	serverAddr := config.Get("SERVER_ADDR", "")
	if serverAddr == "" {
		host := config.Get("SERVER_HOST", "0.0.0.0")
		port := config.Get("SERVER_PORT", "8080")
		serverAddr = host + ":" + port
	}
	server := &http.Server{
		Addr:    serverAddr,
		Handler: httpRouter,
	}

	go func() {
		appLogger.Info("server_starting",
			zap.String("addr", serverAddr),
		)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Fatal("server_error",
				zap.Error(err),
			)
		}
	}()

	<-ctx.Done()
	appLogger.Info("shutdown_signal_received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*1e9)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		appLogger.Error("graceful_shutdown_failed",
			zap.Error(err),
		)
	}
}
