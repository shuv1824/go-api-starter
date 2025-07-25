package core

import (
	"context"

	"github.com/google/uuid"
)

type Service interface {
	Register(ctx context.Context, req CreateUserRequest) (*AuthResponse, error)
	Login(ctx context.Context, req LoginRequest) (*AuthResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
}
