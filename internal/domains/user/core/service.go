package core

import "context"

type Service interface {
	Register(ctx context.Context, req CreateUserRequest) (*AuthResponse, error)
	Login(ctx context.Context, req LoginRequest) (*AuthResponse, error)
}
