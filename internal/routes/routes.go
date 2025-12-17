package routes

import (
    "github.com/gofiber/fiber/v2"
    "github.com/adityaK87/go-backend-assignment/internal/handler"
)

func SetupRoutes(app *fiber.App, userHandler *handler.UserHandler) {
    api := app.Group("/")
    
    // User routes
    users := api.Group("/users")
    users.Post("/", userHandler.CreateUser)
    users.Get("/", userHandler.ListUsers)
    users.Get("/:id", userHandler.GetUser)
    users.Put("/:id", userHandler.UpdateUser)
    users.Delete("/:id", userHandler.DeleteUser)
    
    // Health check
    app.Get("/health", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{
            "status": "ok",
        })
    })
}