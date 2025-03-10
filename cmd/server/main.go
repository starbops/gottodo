package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/starbops/gottodo/internal/handlers"
	"github.com/starbops/gottodo/internal/repositories"
	"github.com/starbops/gottodo/internal/services"
	"github.com/starbops/gottodo/pkg/auth"
	"github.com/starbops/gottodo/pkg/config"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found")
	}

	// Load configuration
	configPath := getConfigPath()
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("Using repository type: %s", cfg.Repository.Type)

	// Create a new Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Initialize repositories
	todoRepo, err := repositories.NewTodoRepository(cfg)
	if err != nil {
		log.Fatalf("Failed to create repository: %v", err)
	}

	// Initialize services
	todoService := services.NewTodoService(todoRepo)

	// Initialize auth service
	authService := auth.NewAuthService()

	// Initialize handlers
	todoHandler := handlers.NewTodoHandler(todoService)
	pageHandler := handlers.NewPageHandler(todoService, authService)
	authHandler := handlers.NewAuthHandler(authService)

	// Auth middleware
	authMiddleware := authHandler.AuthMiddleware

	// Routes
	// Public routes
	e.GET("/", pageHandler.Home)
	e.GET("/login", pageHandler.Login)
	e.GET("/register", pageHandler.Register)

	// Auth routes
	e.GET("/auth/github", authHandler.GitHubAuth)
	e.GET("/auth/github/callback", authHandler.GitHubCallback)
	e.POST("/auth/login", authHandler.Login)
	e.POST("/auth/register", authHandler.Register)
	e.POST("/auth/logout", authHandler.Logout)

	// Protected routes
	e.GET("/dashboard", pageHandler.Dashboard, authMiddleware)

	// Todo API routes
	todoGroup := e.Group("/todos", authMiddleware)
	todoGroup.GET("", todoHandler.GetAllTodos)
	todoGroup.GET("/:id", todoHandler.GetTodo)
	todoGroup.POST("", todoHandler.CreateTodo)
	todoGroup.PUT("/:id", todoHandler.UpdateTodo)
	todoGroup.PUT("/:id/complete", todoHandler.UpdateTodoStatus)
	todoGroup.PUT("/:id/incomplete", todoHandler.UpdateTodoStatus)
	todoGroup.DELETE("/:id", todoHandler.DeleteTodo)

	// Start the server
	port := cfg.Server.Port
	if envPort := os.Getenv("PORT"); envPort != "" {
		// Environment variable overrides config file
		port = envPort
	}

	log.Printf("Server starting on http://localhost:%s", port)
	log.Fatal(e.Start(":" + port))
}

// getConfigPath returns the path to the configuration file
func getConfigPath() string {
	// Check if CONFIG_PATH environment variable is set
	if path := os.Getenv("CONFIG_PATH"); path != "" {
		return path
	}

	// Default to config.json in the current directory
	exePath, err := os.Executable()
	if err != nil {
		log.Println("Warning: Could not determine executable path, using current directory")
		return "config.json"
	}

	return filepath.Join(filepath.Dir(exePath), "config.json")
}
