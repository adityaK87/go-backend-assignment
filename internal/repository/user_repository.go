package repository

import (
    "context"
    "database/sql"
    "time"
    
    "github.com/adityaK87/go-backend-assignment/db/sqlc/generated"
)

type UserRepository interface {
    Create(ctx context.Context, name string, dob time.Time) (*db.User, error)
    GetByID(ctx context.Context, id int32) (*db.User, error)
    List(ctx context.Context, limit, offset int32) ([]*db.User, error)
    Update(ctx context.Context, id int32, name string, dob time.Time) (*db.User, error)
    Delete(ctx context.Context, id int32) error
    Count(ctx context.Context) (int64, error)
}

type userRepository struct {
    queries *db.Queries
}

func NewUserRepository(database *sql.DB) UserRepository {
    return &userRepository{
        queries: db.New(database),
    }
}

func (r *userRepository) Create(ctx context.Context, name string, dob time.Time) (*db.User, error) {
    user, err := r.queries.CreateUser(ctx, db.CreateUserParams{
        Name: name,
        Dob:  dob,
    })
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *userRepository) GetByID(ctx context.Context, id int32) (*db.User, error) {
    user, err := r.queries.GetUserByID(ctx, id)
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *userRepository) List(ctx context.Context, limit, offset int32) ([]*db.User, error) {
    users, err := r.queries.ListUsers(ctx, db.ListUsersParams{
        Limit:  limit,
        Offset: offset,
    })
    if err != nil {
        return nil, err
    }
    
    result := make([]*db.User, len(users))
    for i := range users {
        result[i] = &users[i]
    }
    return result, nil
}

func (r *userRepository) Update(ctx context.Context, id int32, name string, dob time.Time) (*db.User, error) {
    user, err := r.queries.UpdateUser(ctx, db.UpdateUserParams{
        ID:   id,
        Name: name,
        Dob:  dob,
    })
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *userRepository) Delete(ctx context.Context, id int32) error {
    return r.queries.DeleteUser(ctx, id)
}

func (r *userRepository) Count(ctx context.Context) (int64, error) {
    return r.queries.CountUsers(ctx)
}