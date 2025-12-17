package service

import (
    "context"
    "database/sql"
    "errors"
    "time"
    
    "github.com/adityaK87/go-backend-assignment/internal/models"
    "github.com/adityaK87/go-backend-assignment/internal/repository"
    "go.uber.org/zap"
)

type UserService interface {
    CreateUser(ctx context.Context, req models.CreateUserRequest) (*models.UserResponse, error)
    GetUserByID(ctx context.Context, id int32) (*models.UserResponse, error)
    ListUsers(ctx context.Context, page, limit int) ([]*models.UserResponse, error)
    UpdateUser(ctx context.Context, id int32, req models.UpdateUserRequest) (*models.UserResponse, error)
    DeleteUser(ctx context.Context, id int32) error
}

type userService struct {
    repo   repository.UserRepository
    logger *zap.Logger
}

func NewUserService(repo repository.UserRepository, logger *zap.Logger) UserService {
    return &userService{
        repo:   repo,
        logger: logger,
    }
}

func (s *userService) CreateUser(ctx context.Context, req models.CreateUserRequest) (*models.UserResponse, error) {
    // Parse DOB
    dob, err := time.Parse("2006-01-02", req.DOB)
    if err != nil {
        s.logger.Error("Failed to parse DOB", zap.Error(err))
        return nil, errors.New("invalid date format")
    }
    
    // Validate DOB is not in the future
    if dob.After(time.Now()) {
        return nil, errors.New("date of birth cannot be in the future")
    }
    
    // Create user
    user, err := s.repo.Create(ctx, req.Name, dob)
    if err != nil {
        s.logger.Error("Failed to create user", zap.Error(err))
        return nil, err
    }
    
    s.logger.Info("User created successfully", zap.Int32("user_id", user.ID))
    
    return &models.UserResponse{
        ID:   user.ID,
        Name: user.Name,
        DOB:  user.Dob.Format("2006-01-02"),
    }, nil
}

func (s *userService) GetUserByID(ctx context.Context, id int32) (*models.UserResponse, error) {
    user, err := s.repo.GetByID(ctx, id)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, errors.New("user not found")
        }
        s.logger.Error("Failed to get user", zap.Error(err), zap.Int32("user_id", id))
        return nil, err
    }
    
    age := models.CalculateAge(user.Dob)
    
    return &models.UserResponse{
        ID:   user.ID,
        Name: user.Name,
        DOB:  user.Dob.Format("2006-01-02"),
        Age:  &age,
    }, nil
}

func (s *userService) ListUsers(ctx context.Context, page, limit int) ([]*models.UserResponse, error) {
    if page < 1 {
        page = 1
    }
    if limit < 1 || limit > 100 {
        limit = 10
    }
    
    offset := (page - 1) * limit
    
    users, err := s.repo.List(ctx, int32(limit), int32(offset))
    if err != nil {
        s.logger.Error("Failed to list users", zap.Error(err))
        return nil, err
    }
    
    response := make([]*models.UserResponse, len(users))
    for i, user := range users {
        age := models.CalculateAge(user.Dob)
        response[i] = &models.UserResponse{
            ID:   user.ID,
            Name: user.Name,
            DOB:  user.Dob.Format("2006-01-02"),
            Age:  &age,
        }
    }
    
    return response, nil
}

func (s *userService) UpdateUser(ctx context.Context, id int32, req models.UpdateUserRequest) (*models.UserResponse, error) {
    // Check if user exists
    _, err := s.repo.GetByID(ctx, id)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, errors.New("user not found")
        }
        return nil, err
    }
    
    // Parse DOB
    dob, err := time.Parse("2006-01-02", req.DOB)
    if err != nil {
        return nil, errors.New("invalid date format")
    }
    
    if dob.After(time.Now()) {
        return nil, errors.New("date of birth cannot be in the future")
    }
    
    // Update user
    user, err := s.repo.Update(ctx, id, req.Name, dob)
    if err != nil {
        s.logger.Error("Failed to update user", zap.Error(err), zap.Int32("user_id", id))
        return nil, err
    }
    
    s.logger.Info("User updated successfully", zap.Int32("user_id", user.ID))
    
    return &models.UserResponse{
        ID:   user.ID,
        Name: user.Name,
        DOB:  user.Dob.Format("2006-01-02"),
    }, nil
}

func (s *userService) DeleteUser(ctx context.Context, id int32) error {
    // Check if user exists
    _, err := s.repo.GetByID(ctx, id)
    if err != nil {
        if err == sql.ErrNoRows {
            return errors.New("user not found")
        }
        return err
    }
    
    if err := s.repo.Delete(ctx, id); err != nil {
        s.logger.Error("Failed to delete user", zap.Error(err), zap.Int32("user_id", id))
        return err
    }
    
    s.logger.Info("User deleted successfully", zap.Int32("user_id", id))
    return nil
}