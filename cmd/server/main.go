package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/starbops/gottodo/internal/handlers"
	"github.com/starbops/gottodo/internal/repositories"
	"github.com/starbops/gottodo/internal/services"
	"github.com/starbops/gottodo/pkg/auth"
	"github.com/starbops/gottodo/pkg/config"
)

func main() {
	// Define command-line flags
	configPath := flag.String("config", "", "Path to configuration file (default: config.json in executable directory)")
	flag.Parse()

	// If config path is not specified via flag, use the default
	if *configPath == "" {
		*configPath = getDefaultConfigPath()
	}

	// Load configuration
	log.Printf("Loading configuration from %s", *configPath)
	cfg, err := config.LoadConfig(*configPath)
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
	authService := auth.NewAuthService(cfg)

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
	log.Printf("Server starting on http://localhost:%s", port)
	log.Fatal(e.Start(":" + port))
}

// getDefaultConfigPath returns the path to the default configuration file
func getDefaultConfigPath() string {
	// Try to determine the executable path
	exePath, err := os.Executable()
	if err != nil {
		log.Println("Warning: Could not determine executable path, using current directory")
		return "config.json"
	}

	return filepath.Join(filepath.Dir(exePath), "config.json")
}
