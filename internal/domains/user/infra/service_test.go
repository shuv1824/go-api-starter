package infra

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shuv1824/go-api-starter/internal/common/auth"
	apperrors "github.com/shuv1824/go-api-starter/internal/common/errors"
	"github.com/shuv1824/go-api-starter/internal/config"
	"github.com/shuv1824/go-api-starter/internal/domains/user/core"
	"golang.org/x/crypto/bcrypt"
)

// MockUserRepository implements core.UserRepository for testing
type MockUserRepository struct {
	users map[string]*core.User // key: email
	err   error
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: make(map[string]*core.User),
	}
}

func (m *MockUserRepository) Create(ctx context.Context, user *core.User) error {
	if m.err != nil {
		return m.err
	}
	m.users[user.Email] = user
	return nil
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*core.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	for _, user := range m.users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, apperrors.ErrNotFound
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*core.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	user, exists := m.users[email]
	if !exists {
		return nil, apperrors.ErrNotFound
	}
	return user, nil
}

func (m *MockUserRepository) Update(ctx context.Context, user *core.User) error {
	if m.err != nil {
		return m.err
	}
	m.users[user.Email] = user
	return nil
}

func (m *MockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if m.err != nil {
		return m.err
	}
	for email, user := range m.users {
		if user.ID == id {
			delete(m.users, email)
			return nil
		}
	}
	return apperrors.ErrNotFound
}

func (m *MockUserRepository) List(ctx context.Context, limit, offset int) ([]*core.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	var users []*core.User
	for _, user := range m.users {
		users = append(users, user)
	}
	return users, nil
}

func (m *MockUserRepository) Count(ctx context.Context) (int64, error) {
	if m.err != nil {
		return 0, m.err
	}
	return int64(len(m.users)), nil
}

func (m *MockUserRepository) SetError(err error) {
	m.err = err
}

func (m *MockUserRepository) AddUser(user *core.User) {
	m.users[user.Email] = user
}

func setupTestService(t *testing.T) (*service, *MockUserRepository, *auth.Service) {
	// Load test configuration
	cfg, err := config.InitConfig("../../../../config.test.yaml")
	if err != nil {
		t.Skipf("Skipping test: could not load test config: %v", err)
	}

	mockRepo := NewMockUserRepository()
	jwtService := auth.NewService(cfg.Secret, time.Hour, time.Hour*24)
	service := NewService(mockRepo, jwtService)

	return service, mockRepo, jwtService
}

func TestService_Register(t *testing.T) {
	service, mockRepo, _ := setupTestService(t)

	tests := []struct {
		name        string
		request     core.CreateUserRequest
		setupMock   func()
		expectError bool
		errorType   error
	}{
		{
			name: "successful registration",
			request: core.CreateUserRequest{
				Email:    "test@example.com",
				Password: "password123",
				Name:     "Test User",
			},
			setupMock:   func() { mockRepo.SetError(nil) },
			expectError: false,
		},
		{
			name: "email already exists",
			request: core.CreateUserRequest{
				Email:    "existing@example.com",
				Password: "password123",
				Name:     "Test User",
			},
			setupMock: func() {
				mockRepo.SetError(nil)
				existingUser := &core.User{
					ID:       uuid.New(),
					Email:    "existing@example.com",
					Password: "hashedpass",
					Name:     "Existing User",
					IsActive: true,
				}
				mockRepo.AddUser(existingUser)
			},
			expectError: true,
			errorType:   apperrors.ErrEmailExists,
		},
		{
			name: "repository error on GetByEmail",
			request: core.CreateUserRequest{
				Email:    "test@example.com",
				Password: "password123",
				Name:     "Test User",
			},
			setupMock: func() {
				mockRepo.SetError(errors.New("database error"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// service, mockRepo, _ := setupTestService(t)
			tt.setupMock()

			resp, err := service.Register(context.Background(), tt.request)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}
				if tt.errorType != nil && !errors.Is(err, tt.errorType) {
					t.Errorf("expected error %v, got %v", tt.errorType, err)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if resp == nil {
				t.Error("expected response but got nil")
				return
			}

			if resp.User.Email != tt.request.Email {
				t.Errorf("expected email %s, got %s", tt.request.Email, resp.User.Email)
			}

			if resp.User.Name != tt.request.Name {
				t.Errorf("expected name %s, got %s", tt.request.Name, resp.User.Name)
			}

			if resp.Token == "" {
				t.Error("expected token but got empty string")
			}

			// Verify password was hashed
			if resp.User.Password == tt.request.Password {
				t.Error("password should be hashed")
			}
		})
	}
}

func TestService_Login(t *testing.T) {
	service, mockRepo, _ := setupTestService(t)

	// Create a test user with hashed password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	testUser := &core.User{
		ID:       uuid.New(),
		Email:    "test@example.com",
		Password: string(hashedPassword),
		Name:     "Test User",
		IsActive: true,
	}

	tests := []struct {
		name        string
		request     core.LoginRequest
		setupMock   func()
		expectError bool
		errorType   error
	}{
		{
			name: "successful login",
			request: core.LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMock: func() {
				mockRepo.SetError(nil)
				mockRepo.AddUser(testUser)
			},
			expectError: false,
		},
		{
			name: "user not found",
			request: core.LoginRequest{
				Email:    "nonexistent@example.com",
				Password: "password123",
			},
			setupMock: func() {
				mockRepo.SetError(nil)
			},
			expectError: true,
			errorType:   apperrors.ErrInvalidPassword,
		},
		{
			name: "invalid password",
			request: core.LoginRequest{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			setupMock: func() {
				mockRepo.SetError(nil)
				mockRepo.AddUser(testUser)
			},
			expectError: true,
			errorType:   apperrors.ErrInvalidPassword,
		},
		{
			name: "inactive user",
			request: core.LoginRequest{
				Email:    "inactive@example.com",
				Password: "password123",
			},
			setupMock: func() {
				mockRepo.SetError(nil)
				inactiveUser := &core.User{
					ID:       uuid.New(),
					Email:    "inactive@example.com",
					Password: string(hashedPassword),
					Name:     "Inactive User",
					IsActive: false,
				}
				mockRepo.AddUser(inactiveUser)
			},
			expectError: true,
			errorType:   apperrors.ErrUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// service, mockRepo, _ := setupTestService(t)
			tt.setupMock()

			resp, err := service.Login(context.Background(), tt.request)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}
				if tt.errorType != nil && !errors.Is(err, tt.errorType) {
					t.Errorf("expected error %v, got %v", tt.errorType, err)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if resp == nil {
				t.Error("expected response but got nil")
				return
			}

			if resp.User.Email != tt.request.Email {
				t.Errorf("expected email %s, got %s", tt.request.Email, resp.User.Email)
			}

			if resp.Token == "" {
				t.Error("expected token but got empty string")
			}
		})
	}
}
