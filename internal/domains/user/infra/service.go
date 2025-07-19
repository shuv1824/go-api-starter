package infra

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/shuv1824/go-api-starter/internal/common/auth"
	"github.com/shuv1824/go-api-starter/internal/domains/user/core"
	"golang.org/x/crypto/bcrypt"

	apperrors "github.com/shuv1824/go-api-starter/internal/common/errors"
)

type service struct {
	repo       core.UserRepository
	jwtService *auth.Service
}

func NewService(repo core.UserRepository, jwtService *auth.Service) *service {
	return &service{
		repo:       repo,
		jwtService: jwtService,
	}
}

func (s *service) Register(ctx context.Context, req core.CreateUserRequest) (*core.AuthResponse, error) {
	// Check if user already exists
	existingUser, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, apperrors.ErrNotFound) {
		return nil, err
	}
	if existingUser != nil {
		return nil, apperrors.ErrEmailExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	id := uuid.New()

	// Create user
	user := &core.User{
		ID:       id,
		Email:    req.Email,
		Password: string(hashedPassword),
		Name:     req.Name,
		IsActive: true,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Generate JWT token
	token, err := s.jwtService.GenerateToken(user.ID.String(), user.Email)
	if err != nil {
		return nil, err
	}

	return &core.AuthResponse{
		Token: token,
		User:  *user,
	}, nil
}

func (s *service) Login(ctx context.Context, req core.LoginRequest) (*core.AuthResponse, error) {
	user, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return nil, apperrors.ErrInvalidPassword
		}
		return nil, err
	}

	if !user.IsActive {
		return nil, apperrors.ErrUnauthorized
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, apperrors.ErrInvalidPassword
	}

	// Generate JWT token
	token, err := s.jwtService.GenerateToken(user.ID.String(), user.Email)
	if err != nil {
		return nil, err
	}

	return &core.AuthResponse{
		Token: token,
		User:  *user,
	}, nil
}
