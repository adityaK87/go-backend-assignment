package middleware

import (
    "time"
    
    "github.com/gofiber/fiber/v2"
    "go.uber.org/zap"
)

func Logger(logger *zap.Logger) fiber.Handler {
    return func(c *fiber.Ctx) error {
        start := time.Now()
        
        // Process request
        err := c.Next()
        
        // Calculate duration
        duration := time.Since(start)
        
        // Get request ID
        requestID, _ := c.Locals("requestID").(string)
        
        // Log request
        logger.Info("HTTP Request",
            zap.String("request_id", requestID),
            zap.String("method", c.Method()),
            zap.String("path", c.Path()),
            zap.Int("status", c.Response().StatusCode()),
            zap.Duration("duration", duration),
            zap.String("ip", c.IP()),
        )
        
        return err
    }
}