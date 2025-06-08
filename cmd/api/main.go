package main

import (
	"database/sql"
	"fmt"
	"github.com/evgeniySeleznev/person-enrichment-service/internal/delivery/http"
	"github.com/evgeniySeleznev/person-enrichment-service/internal/repository/api"
	"github.com/evgeniySeleznev/person-enrichment-service/internal/repository/postgresql"
	"github.com/evgeniySeleznev/person-enrichment-service/internal/server"
	"github.com/evgeniySeleznev/person-enrichment-service/internal/service"
	"github.com/evgeniySeleznev/person-enrichment-service/pkg/logger"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"os"

	"context"
	"time"

	_ "github.com/lib/pq"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// @title Person Enrichment API
// @version 1.0
// @description Сервис для обогащения данных о людях

// @host localhost:8080
// @BasePath /api
func main() {
	appLogger := logger.NewLogger(os.Stdout)
	appLogger.Info("Starting application...")

	db, err := initDB(appLogger)
	if err != nil {
		appLogger.Fatal("Database initialization failed", err)
	}
	defer db.Close()

	personService := initServices(db, appLogger)
	router := http.NewRouter(personService, appLogger)

	server := server.NewServer(os.Getenv("APP_Port"), router, appLogger)
	server.Start()

	//server := NewServer(os.Getenv("APP_PORT"), router, appLogger)
	//server.Start()
}

func initDB(logger logger.Logger) (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSL_MODE"),
	))
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Проверка подключения (уже было)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}
	logger.Info("Database connection established")

	// Добавляем автоматические миграции (новый код)
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to create migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to init migrations: %w", err)
	}

	logger.Info("Applying database migrations...")
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, fmt.Errorf("failed to apply migrations: %w", err)
	}
	logger.Info("Migrations applied successfully")

	return db, nil
}

func initServices(db *sql.DB, logger logger.Logger) *service.PersonService {
	personRepo := postgresql.NewPersonRepository(db)
	apiClient := api.NewAPIClient(
		os.Getenv("AGIFY_URL"),
		os.Getenv("GENDERIZE_URL"),
		os.Getenv("NATIONALIZE_URL"),
	)
	return service.NewPersonService(personRepo, apiClient)
}
