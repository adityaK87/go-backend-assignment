package handler

import (
    "strconv"
    
    "github.com/go-playground/validator/v10"
    "github.com/gofiber/fiber/v2"
    "github.com/adityaK87/go-backend-assignment/internal/models"
    "github.com/adityaK87/go-backend-assignment/internal/service"
    "go.uber.org/zap"
)

type UserHandler struct {
    service   service.UserService
    validator *validator.Validate
    logger    *zap.Logger
}

func NewUserHandler(service service.UserService, logger *zap.Logger) *UserHandler {
    return &UserHandler{
        service:   service,
        validator: validator.New(),
        logger:    logger,
    }
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
    var req models.CreateUserRequest
    
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
            Error: "Invalid request body",
        })
    }
    
    if err := h.validator.Struct(req); err != nil {
        details := make(map[string]string)
        for _, err := range err.(validator.ValidationErrors) {
            details[err.Field()] = err.Tag()
        }
        return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
            Error:   "Validation failed",
            Details: details,
        })
    }
    
    user, err := h.service.CreateUser(c.Context(), req)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
            Error: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusCreated).JSON(user)
}

func (h *UserHandler) GetUser(c *fiber.Ctx) error {
    id, err := strconv.ParseInt(c.Params("id"), 10, 32)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
            Error: "Invalid user ID",
        })
    }
    
    user, err := h.service.GetUserByID(c.Context(), int32(id))
    if err != nil {
        if err.Error() == "user not found" {
            return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
                Error: "User not found",
            })
        }
        return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
            Error: err.Error(),
        })
    }
    
    return c.JSON(user)
}

func (h *UserHandler) ListUsers(c *fiber.Ctx) error {
    var pagination models.PaginationQuery
    
    if err := c.QueryParser(&pagination); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
            Error: "Invalid query parameters",
        })
    }
    
    if pagination.Page == 0 {
        pagination.Page = 1
    }
    if pagination.Limit == 0 {
        pagination.Limit = 10
    }
    
    users, err := h.service.ListUsers(c.Context(), pagination.Page, pagination.Limit)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
            Error: err.Error(),
        })
    }
    
    return c.JSON(users)
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
    id, err := strconv.ParseInt(c.Params("id"), 10, 32)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
            Error: "Invalid user ID",
        })
    }
    
    var req models.UpdateUserRequest
    
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
            Error: "Invalid request body",
        })
    }
    
    if err := h.validator.Struct(req); err != nil {
        details := make(map[string]string)
        for _, err := range err.(validator.ValidationErrors) {
            details[err.Field()] = err.Tag()
        }
        return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
            Error:   "Validation failed",
            Details: details,
        })
    }
    
    user, err := h.service.UpdateUser(c.Context(), int32(id), req)
    if err != nil {
        if err.Error() == "user not found" {
            return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
                Error: "User not found",
            })
        }
        return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
            Error: err.Error(),
        })
    }
    
    return c.JSON(user)
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
    id, err := strconv.ParseInt(c.Params("id"), 10, 32)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
            Error: "Invalid user ID",
        })
    }
    
    if err := h.service.DeleteUser(c.Context(), int32(id)); err != nil {
        if err.Error() == "user not found" {
            return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
                Error: "User not found",
            })
        }
        return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
            Error: err.Error(),
        })
    }
    
    return c.SendStatus(fiber.StatusNoContent)
}