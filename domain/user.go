package domain

import "context"

type UserUsecase interface {
	GetUsers(ctx context.Context) ([]*User, error)
	GetUserByIdentifier(ctx context.Context, identifier string, value string) (*User, error)
}

type UserRepository interface {
	GetUsers(ctx context.Context) ([]*User, error)
	GetUserByUsername(ctx context.Context, username string) (*User, error)
}

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
}
