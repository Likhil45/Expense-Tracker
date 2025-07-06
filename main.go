package main

import (
	"context"
	"expense-tracker/auth"
	"expense-tracker/controller"
	_ "expense-tracker/docs"
	"expense-tracker/postgresql"

	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Expense Tracker API
// @version 1.0
// @description This is a sample server for an expense tracker.
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	// Load .env file
	_ = godotenv.Load()

	id := uuid.New()
	println(id.String())

	postgresql.ConnectPostgres()

	s := gin.Default()

	// Public routes
	s.POST("/api/v1/login", controller.Login)
	s.POST("/api/v1/signup", controller.SignUp)

	// Protected routes with JWT and Rate Limiting
	r := s.Group("/api/v1/expenses")
	r.Use(auth.JWTAuthMiddleware(), auth.RateLimitMiddleware())
	r.POST("/", controller.CreateExpense)
	r.GET("/:id", controller.GetExpenseById)
	r.PUT("/:id", controller.UpdateExpense)
	r.DELETE("/:id", controller.DeleteExpense)
	r.GET("/", controller.ListExpensesWithFilters)
	r.GET("/summary", controller.Summary)

	s.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	srv := &http.Server{
		Addr:    ":8080",
		Handler: s,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Info("Server exiting")
}
