package infra

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/shuv1824/go-api-starter/internal/config"
	"github.com/shuv1824/go-api-starter/internal/migration"
	"github.com/shuv1824/go-api-starter/pkg/database"
	apperrors "github.com/shuv1824/go-api-starter/internal/common/errors"
	"github.com/shuv1824/go-api-starter/internal/domains/user/core"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	// Load test configuration
	cfg, err := config.InitConfig("../../../../config.test.yaml")
	if err != nil {
		t.Skipf("Skipping test: could not load test config: %v", err)
	}

	// Create unique test database name to avoid conflicts
	testDBName := fmt.Sprintf("%s_%s", cfg.Database.DbName, uuid.New().String()[:8])

	// First connect to postgres database to create test database
	pgCfg := cfg.Database
	pgCfg.DbName = "postgres" // Connect to default postgres db first

	pgDB, err := database.NewDatabase(&pgCfg)
	if err != nil {
		t.Skipf("Skipping test: could not connect to postgres: %v", err)
	}

	// Create test database
	sqlDB, err := pgDB.DB()
	if err != nil {
		t.Fatalf("failed to get database instance: %v", err)
	}

	_, err = sqlDB.Exec(fmt.Sprintf("CREATE DATABASE %s", testDBName))
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}
	sqlDB.Close()

	// Now connect to the test database
	testCfg := cfg.Database
	testCfg.DbName = testDBName

	db, err := database.NewDatabase(&testCfg)
	if err != nil {
		t.Skipf("Skipping test: could not connect to test database: %v", err)
	}

	// Run migrations
	if err := migration.MigrateUp(db, cfg.Database.Type); err != nil {
		t.Fatalf("failed to migrate test database: %v", err)
	}

	// Register cleanup function
	t.Cleanup(func() {
		// Close current connection
		sqlDB, err := db.DB()
		if err == nil {
			sqlDB.Close()
		}

		// Connect to postgres db to drop test database
		pgDB, err := database.NewDatabase(&pgCfg)
		if err == nil {
			sqlDB, err := pgDB.DB()
			if err == nil {
				sqlDB.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", testDBName))
				sqlDB.Close()
			}
		}
	})

	return db
}

func TestMain(m *testing.M) {
	// Check if PostgreSQL is available for testing
	cfg, err := config.InitConfig("../../../../config.test.yaml")
	if err != nil {
		fmt.Printf("Warning: Test config not found, skipping database tests: %v\n", err)
		os.Exit(0)
	}

	// Try to connect to test database to ensure it's available
	pgCfg := cfg.Database
	pgCfg.DbName = "postgres" // Connect to default postgres DB first
	
	_, err = database.NewDatabase(&pgCfg)
	if err != nil {
		fmt.Printf("Warning: PostgreSQL not available for testing, skipping database tests: %v\n", err)
		os.Exit(0)
	}

	// Run tests
	os.Exit(m.Run())
}

func TestUserRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)

	user := &core.User{
		ID:       uuid.New(),
		Email:    "test@example.com",
		Password: "hashedpassword",
		Name:     "Test User",
		IsActive: true,
	}

	err := repo.Create(context.Background(), user)
	if err != nil {
		t.Errorf("unexpected error creating user: %v", err)
	}

	// Verify user was created
	var dbUser core.User
	result := db.Where("email = ?", user.Email).First(&dbUser)
	if result.Error != nil {
		t.Errorf("user not found in database: %v", result.Error)
	}

	if dbUser.Email != user.Email {
		t.Errorf("expected email %s, got %s", user.Email, dbUser.Email)
	}
}

func TestUserRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)

	// Create a test user
	user := &core.User{
		ID:       uuid.New(),
		Email:    "test@example.com",
		Password: "hashedpassword",
		Name:     "Test User",
		IsActive: true,
	}

	if err := db.Create(user).Error; err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	tests := []struct {
		name        string
		userID      uuid.UUID
		expectError bool
		errorType   error
	}{
		{
			name:        "existing user",
			userID:      user.ID,
			expectError: false,
		},
		{
			name:        "non-existing user",
			userID:      uuid.New(),
			expectError: true,
			errorType:   apperrors.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.GetByID(context.Background(), tt.userID)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}
				if tt.errorType != nil && err != tt.errorType {
					t.Errorf("expected error %v, got %v", tt.errorType, err)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if result == nil {
				t.Error("expected user but got nil")
				return
			}

			if result.ID != user.ID {
				t.Errorf("expected user ID %s, got %s", user.ID, result.ID)
			}
		})
	}
}

func TestUserRepository_GetByEmail(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)

	// Create a test user
	user := &core.User{
		ID:       uuid.New(),
		Email:    "test@example.com",
		Password: "hashedpassword",
		Name:     "Test User",
		IsActive: true,
	}

	if err := db.Create(user).Error; err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	tests := []struct {
		name        string
		email       string
		expectError bool
		errorType   error
	}{
		{
			name:        "existing user",
			email:       user.Email,
			expectError: false,
		},
		{
			name:        "non-existing user",
			email:       "nonexistent@example.com",
			expectError: true,
			errorType:   apperrors.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.GetByEmail(context.Background(), tt.email)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}
				if tt.errorType != nil && err != tt.errorType {
					t.Errorf("expected error %v, got %v", tt.errorType, err)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if result == nil {
				t.Error("expected user but got nil")
				return
			}

			if result.Email != user.Email {
				t.Errorf("expected email %s, got %s", user.Email, result.Email)
			}
		})
	}
}

func TestUserRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)

	// Create a test user
	user := &core.User{
		ID:       uuid.New(),
		Email:    "test@example.com",
		Password: "hashedpassword",
		Name:     "Test User",
		IsActive: true,
	}

	if err := db.Create(user).Error; err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	// Update the user
	user.Name = "Updated User"
	user.IsActive = false

	err := repo.Update(context.Background(), user)
	if err != nil {
		t.Errorf("unexpected error updating user: %v", err)
	}

	// Verify user was updated
	var dbUser core.User
	result := db.Where("id = ?", user.ID).First(&dbUser)
	if result.Error != nil {
		t.Errorf("user not found in database: %v", result.Error)
	}

	if dbUser.Name != "Updated User" {
		t.Errorf("expected name 'Updated User', got %s", dbUser.Name)
	}

	if dbUser.IsActive != false {
		t.Errorf("expected IsActive to be false, got %v", dbUser.IsActive)
	}
}

func TestUserRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)

	// Create a test user
	user := &core.User{
		ID:       uuid.New(),
		Email:    "test@example.com",
		Password: "hashedpassword",
		Name:     "Test User",
		IsActive: true,
	}

	if err := db.Create(user).Error; err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	// Delete the user
	err := repo.Delete(context.Background(), user.ID)
	if err != nil {
		t.Errorf("unexpected error deleting user: %v", err)
	}

	// Verify user was deleted (soft delete)
	var dbUser core.User
	result := db.Where("id = ?", user.ID).First(&dbUser)
	if result.Error != nil {
		// User should not be found due to soft delete
		if result.Error != gorm.ErrRecordNotFound {
			t.Errorf("unexpected error: %v", result.Error)
		}
	}

	// Verify user exists in database with deleted_at set (for soft delete)
	result = db.Unscoped().Where("id = ?", user.ID).First(&dbUser)
	if result.Error != nil {
		t.Errorf("user should still exist in database with deleted_at set: %v", result.Error)
	}
}

func TestUserRepository_List(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)

	// Create test users
	users := []*core.User{
		{
			ID:       uuid.New(),
			Email:    "user1@example.com",
			Password: "hashedpassword",
			Name:     "User 1",
			IsActive: true,
		},
		{
			ID:       uuid.New(),
			Email:    "user2@example.com",
			Password: "hashedpassword",
			Name:     "User 2",
			IsActive: true,
		},
		{
			ID:       uuid.New(),
			Email:    "user3@example.com",
			Password: "hashedpassword",
			Name:     "User 3",
			IsActive: true,
		},
	}

	for _, user := range users {
		if err := db.Create(user).Error; err != nil {
			t.Fatalf("failed to create test user: %v", err)
		}
	}

	tests := []struct {
		name          string
		limit         int
		offset        int
		expectedCount int
	}{
		{
			name:          "get all users",
			limit:         10,
			offset:        0,
			expectedCount: 3,
		},
		{
			name:          "get first 2 users",
			limit:         2,
			offset:        0,
			expectedCount: 2,
		},
		{
			name:          "get users with offset",
			limit:         10,
			offset:        1,
			expectedCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.List(context.Background(), tt.limit, tt.offset)

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if len(result) != tt.expectedCount {
				t.Errorf("expected %d users, got %d", tt.expectedCount, len(result))
			}
		})
	}
}

func TestUserRepository_Count(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)

	// Initially should be 0
	count, err := repo.Count(context.Background())
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if count != 0 {
		t.Errorf("expected count 0, got %d", count)
	}

	// Create test users
	users := []*core.User{
		{
			ID:       uuid.New(),
			Email:    "user1@example.com",
			Password: "hashedpassword",
			Name:     "User 1",
			IsActive: true,
		},
		{
			ID:       uuid.New(),
			Email:    "user2@example.com",
			Password: "hashedpassword",
			Name:     "User 2",
			IsActive: true,
		},
	}

	for _, user := range users {
		if err := db.Create(user).Error; err != nil {
			t.Fatalf("failed to create test user: %v", err)
		}
	}

	count, err = repo.Count(context.Background())
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if count != 2 {
		t.Errorf("expected count 2, got %d", count)
	}
}
