package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"github.com/swaggo/http-swagger"

	"github.com/evgeniySeleznev/person-enrichment-service/internal/delivery/http"
	"github.com/evgeniySeleznev/person-enrichment-service/internal/repository/api"
	"github.com/evgeniySeleznev/person-enrichment-service/internal/repository/postgresql"
	"github.com/evgeniySeleznev/person-enrichment-service/internal/service"
	"github.com/evgeniySeleznev/person-enrichment-service/pkg/logger"
)

// @title Person Enrichment API
// @version 1.0
// @description Сервис для обогащения данных о людях с использованием внешних API (agify.io, genderize.io, nationalize.io)

// @contact.name API Support
// @contact.url https://github.com/evgeniySeleznev/person-enrichment-service
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api
// @schemes http

func main() {
	// 1. Инициализация логгера
	appLogger := logger.NewLogger(os.Stdout)
	appLogger.Info("Starting application...")

	// 2. Подключение к БД
	db, err := sql.Open("postgres",
		fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
			os.Getenv("DB_SSL_MODE"),
		))
	if err != nil {
		appLogger.Fatal("Failed to connect to database", err)
	}
	defer db.Close()

	// 3. Проверка соединения с БД
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		appLogger.Fatal("Database ping failed", err)
	}

	// 4. Инициализация репозиториев
	personRepo := postgresql.NewPersonRepository(db)
	apiClient := api.NewAPIClient(
		os.Getenv("AGIFY_URL"),
		os.Getenv("GENDERIZE_URL"),
		os.Getenv("NATIONALIZE_URL"),
	)

	// 5. Создание сервисов
	personService := service.NewPersonService(personRepo, apiClient)

	// 6. Настройка HTTP-сервера
	router := http.NewRouter(personService, appLogger)
	server := &http.Server{
		Addr:    ":" + os.Getenv("APP_PORT"),
		Handler: router,
	}

	// 7. Graceful shutdown
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Fatal("Server error", err)
		}
	}()

	appLogger.Info("Server started on port " + os.Getenv("APP_PORT"))

	// Ожидание сигналов завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("Shutting down server...")
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		appLogger.Error("Server shutdown error", err)
	}
}
