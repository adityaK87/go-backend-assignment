package middleware

import (
    "github.com/gofiber/fiber/v2"
    "go.uber.org/zap"
)

func Recover(logger *zap.Logger) fiber.Handler {
    return func(c *fiber.Ctx) error {
        defer func() {
            if r := recover(); r != nil {
                logger.Error("Panic recovered",
                    zap.Any("error", r),
                    zap.String("path", c.Path()),
                )
                
                c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                    "error": "Internal server error",
                })
            }
        }()
        
        return c.Next()
    }
}