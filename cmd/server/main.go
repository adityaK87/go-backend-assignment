package main

import (
    "database/sql"
    "fmt"
    "log"
    "os"
    "os/signal"
    "syscall"
    
    
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
    _ "github.com/lib/pq"
    "go.uber.org/zap"
    
    "github.com/adityaK87/go-backend-assignment/config"
    "github.com/adityaK87/go-backend-assignment/internal/handler"
    "github.com/adityaK87/go-backend-assignment/internal/logger"
    "github.com/adityaK87/go-backend-assignment/internal/middleware"
    "github.com/adityaK87/go-backend-assignment/internal/repository"
    "github.com/adityaK87/go-backend-assignment/internal/routes"
    "github.com/adityaK87/go-backend-assignment/internal/service"
)

func main() {
    // Load configuration
    cfg := config.Load()
    
    // Initialize logger
    if err := logger.Init(); err != nil {
        log.Fatal("Failed to initialize logger:", err)
    }
    defer logger.Sync()
    

    
    // Connect to database
    db, err := sql.Open("postgres", cfg.DatabaseURL)
    if err != nil {
        logger.Log.Fatal("Failed to connect to database", zap.Error(err))
    }
    defer db.Close()
    
    // Verify database connection
    if err := db.Ping(); err != nil {
        logger.Log.Fatal("Failed to ping database", zap.Error(err))
    }
    
    logger.Log.Info("Successfully connected to database")
    
    // Initialize layers
    userRepo := repository.NewUserRepository(db)
    userService := service.NewUserService(userRepo, logger.Log)
    userHandler := handler.NewUserHandler(userService, logger.Log)
    
    // Create Fiber app
    app := fiber.New(fiber.Config{
        ErrorHandler: func(c *fiber.Ctx, err error) error {
            code := fiber.StatusInternalServerError
            if e, ok := err.(*fiber.Error); ok {
                code = e.Code
            }
            return c.Status(code).JSON(fiber.Map{
                "error": err.Error(),
            })
        },
    })
    
    // Middleware
    app.Use(cors.New())
    app.Use(middleware.RequestID())
    app.Use(middleware.Logger(logger.Log))
    app.Use(middleware.Recover(logger.Log))
    
    // Setup routes
    routes.SetupRoutes(app, userHandler)
    
    // Graceful shutdown
    go func() {
        sigChan := make(chan os.Signal, 1)
        signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
        <-sigChan
        
        logger.Log.Info("Shutting down server...")
        app.Shutdown()
    }()
    
    // Start server
    addr := fmt.Sprintf(":%s", cfg.ServerPort)
    logger.Log.Info("Starting server", zap.String("address", addr))
    
    if err := app.Listen(addr); err != nil {
        logger.Log.Fatal("Failed to start server", zap.Error(err))
    }
}