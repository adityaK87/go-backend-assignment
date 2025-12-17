package middleware

import (
    "github.com/gofiber/fiber/v2"
    "github.com/google/uuid"
)

func RequestID() fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Get request ID from header or generate new one
        requestID := c.Get("X-Request-ID")
        if requestID == "" {
            requestID = uuid.New().String()
        }
        
        // Set request ID in context and response header
        c.Locals("requestID", requestID)
        c.Set("X-Request-ID", requestID)
        
        return c.Next()
    }
}