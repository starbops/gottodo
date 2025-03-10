package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/starbops/gottodo/internal/handlers"
	"github.com/starbops/gottodo/internal/repositories"
	"github.com/starbops/gottodo/internal/services"
	"github.com/starbops/gottodo/pkg/auth"
	"github.com/starbops/gottodo/pkg/database"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Initialize database connection
	db, err := database.NewSupabaseClient()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	todoRepo := repositories.NewSupabaseTodoRepository(db)

	// Initialize auth service
	authService := auth.NewAuthService()

	// Initialize services
	todoService := services.NewTodoService(todoRepo)

	// Initialize handlers
	todoHandler := handlers.NewTodoHandler(todoService)
	authHandler := handlers.NewAuthHandler(authService)
	pageHandler := handlers.NewPageHandler(todoService, authService)

	// Initialize Echo server
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Static files
	e.Static("/static", "ui/static")

	// Routes
	// Public routes
	e.GET("/", pageHandler.Home)
	e.GET("/login", pageHandler.Login)
	e.GET("/register", pageHandler.Register)
	e.POST("/auth/register", authHandler.Register)
	e.POST("/auth/login", authHandler.Login)
	e.POST("/auth/logout", authHandler.Logout)

	// GitHub OAuth routes
	e.GET("/auth/github", authHandler.GitHubAuth)
	e.GET("/auth/github/callback", authHandler.GitHubCallback)

	// Protected routes
	dashboardGroup := e.Group("/dashboard")
	dashboardGroup.Use(authHandler.AuthMiddleware)
	dashboardGroup.GET("", pageHandler.Dashboard)

	todoRoutes := e.Group("/todos")
	todoRoutes.Use(authHandler.AuthMiddleware)
	todoRoutes.GET("", todoHandler.GetAllTodos)
	todoRoutes.POST("", todoHandler.CreateTodo)
	todoRoutes.GET("/:id", todoHandler.GetTodo)
	todoRoutes.PUT("/:id", todoHandler.UpdateTodo)
	todoRoutes.DELETE("/:id", todoHandler.DeleteTodo)
	todoRoutes.PUT("/:id/complete", todoHandler.CompleteTodo)
	todoRoutes.PUT("/:id/incomplete", todoHandler.IncompleteTodo)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	e.Logger.Fatal(e.Start(":" + port))
}
