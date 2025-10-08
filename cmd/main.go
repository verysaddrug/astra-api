package main

import (
	"astra-api/internal/cache"
	"astra-api/internal/config"
	"astra-api/internal/handler"
	. "astra-api/internal/handler"
	"astra-api/internal/middleware"
	"astra-api/internal/repository"
	"astra-api/internal/service"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Astra API
// @version 1.0
// @description REST API для хранения и раздачи документов с кэшированием
// @contact.name мой email
// @contact.email verysaddrug@icloud.com
// @host localhost:8080
// @BasePath /
// @schemes http
func main() {
	cfg := initConfig()
	db := initDB(cfg, 10, 3*time.Second)
	if cfg.AutoMigrate {
		migrate(db)
	} else {
		log.Println("Auto-migration disabled. Use AUTO_MIGRATE=true to enable automatic migrations.")
	}

	// Initialize repositories (implementing interfaces)
	var userRepo repository.UserRepositoryInterface = repository.NewUserRepository(db)
	var docRepo repository.DocumentRepositoryInterface = repository.NewDocumentRepository(db)

	// Initialize services (implementing interfaces)
	var authService service.AuthServiceInterface = service.NewAuthService(userRepo, cfg.AdminToken)
	var sessionService service.SessionServiceInterface = service.NewSessionService()
	var docsService service.DocsServiceInterface = service.NewDocsService(docRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService, sessionService)
	cache := cache.NewCache(5 * time.Minute)
	docsHandler := handler.NewDocsHandler(docsService, cache, sessionService, userRepo)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(sessionService, userRepo)

	routes(authHandler, docsHandler, authMiddleware)
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func initConfig() *config.Config {
	cfg := config.LoadConfig(".env")
	if cfg.DBHost == "" || cfg.DBUser == "" || cfg.DBName == "" {
		log.Fatal("DB config is not set properly")
	}
	return cfg
}

func initDB(cfg *config.Config, attempts int, delay time.Duration) *sqlx.DB {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	var db *sqlx.DB
	var err error
	for i := 0; i < attempts; i++ {
		db, err = sqlx.Connect("postgres", dsn)
		if err == nil {
			return db
		}
		log.Printf("DB connection failed (attempt %d/%d): %v", i+1, attempts, err)
		time.Sleep(delay)
	}
	log.Fatalf("Could not connect to DB after %d attempts: %v", attempts, err)
	return nil
}

func migrate(db *sqlx.DB) {
	migrationsDir := "./migrations"
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		log.Printf("Migrations dir %s not found, skipping migrations", migrationsDir)
		return
	}
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("Goose set dialect error: %v", err)
	}
	if err := goose.Up(db.DB, migrationsDir); err != nil {
		log.Fatalf("Goose migration error: %v", err)
	}
	log.Println("Migrations applied successfully")
}

func routes(authHandler *handler.AuthHandler, docsHandler *handler.DocsHandler, authMiddleware *middleware.AuthMiddleware) {
	// Base middleware for all routes
	baseMiddleware := middleware.ChainMiddleware(
		middleware.LoggingMiddleware,
	)

	// Public routes (no auth required)
	http.HandleFunc("/api/register", baseMiddleware(authHandler.Register))
	http.HandleFunc("/api/auth", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			baseMiddleware(authHandler.Auth)(w, r)
		case http.MethodDelete:
			baseMiddleware(authHandler.Logout)(w, r)
		default:
			WriteError(w, 405, "method not allowed")
		}
	})

	http.HandleFunc("/api/auth/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodDelete:
			baseMiddleware(authHandler.Logout)(w, r)
		default:
			WriteError(w, 405, "method not allowed")
		}
	})

	// Protected routes (auth required)
	protectedMiddleware := middleware.ChainMiddleware(
		middleware.LoggingMiddleware,
		authMiddleware.RequireAuth,
	)

	http.HandleFunc("/api/docs", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			protectedMiddleware(docsHandler.Upload)(w, r)
		case http.MethodGet, http.MethodHead:
			protectedMiddleware(docsHandler.List)(w, r)
		default:
			WriteError(w, 405, "method not allowed")
		}
	})

	http.HandleFunc("/api/docs/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet, http.MethodHead:
			protectedMiddleware(docsHandler.GetByID)(w, r)
		case http.MethodDelete:
			protectedMiddleware(docsHandler.DeleteByID)(w, r)
		default:
			WriteError(w, 405, "method not allowed")
		}
	})

	http.Handle("/docs/", http.StripPrefix("/docs/", http.FileServer(http.Dir("./docs"))))
	httpSwagger.URL("http://localhost:8080/docs/swagger.json")
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)
	http.HandleFunc("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "./docs/swagger.json") })
}
