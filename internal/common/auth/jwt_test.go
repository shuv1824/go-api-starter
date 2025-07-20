package auth

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/shuv1824/go-api-starter/internal/config"
)

func setupTestAuthService(t *testing.T) *Service {
	// Load test configuration
	cfg, err := config.InitConfig("../../../config.test.yaml")
	if err != nil {
		// Fallback to default values if config not available
		return NewService("test-secret-key", time.Hour, time.Hour*24)
	}
	return NewService(cfg.Secret, time.Hour, time.Hour*24)
}

func TestMain(m *testing.M) {
	// Check if test config is available
	_, err := config.InitConfig("../../../config.test.yaml")
	if err != nil {
		fmt.Printf("Warning: Test config not found, using default values: %v\n", err)
	}

	// Run tests
	os.Exit(m.Run())
}

func TestService_GenerateToken(t *testing.T) {
	service := setupTestAuthService(t)
	
	userID := uuid.New().String()
	email := "test@example.com"

	token, err := service.GenerateToken(userID, email)
	if err != nil {
		t.Errorf("unexpected error generating token: %v", err)
	}

	if token == "" {
		t.Error("expected token but got empty string")
	}

	// Verify the token can be parsed
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return service.secretKey, nil
	})

	if err != nil {
		t.Errorf("error parsing generated token: %v", err)
	}

	if !parsedToken.Valid {
		t.Error("generated token is not valid")
	}

	// Verify claims
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok {
		if claims["user_id"] != userID {
			t.Errorf("expected user_id %s, got %v", userID, claims["user_id"])
		}
		if claims["email"] != email {
			t.Errorf("expected email %s, got %v", email, claims["email"])
		}
		if claims["sub"] != userID {
			t.Errorf("expected subject %s, got %v", userID, claims["sub"])
		}
	} else {
		t.Error("could not parse token claims")
	}
}

func TestService_ValidateToken(t *testing.T) {
	service := setupTestAuthService(t)
	userID := uuid.New().String()
	email := "test@example.com"

	// Generate a valid token
	validToken, err := service.GenerateToken(userID, email)
	if err != nil {
		t.Fatalf("failed to generate test token: %v", err)
	}

	// Create an expired token
	expiredClaims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
		},
	}
	expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
	expiredTokenString, _ := expiredToken.SignedString(service.secretKey)

	// Create a token with wrong signature
	wrongSecretService := NewService("wrong-secret", time.Hour, time.Hour*24)
	wrongSignatureToken, _ := wrongSecretService.GenerateToken(userID, email)

	tests := []struct {
		name        string
		token       string
		expectError bool
		expectEmail string
		expectID    string
	}{
		{
			name:        "valid token",
			token:       validToken,
			expectError: false,
			expectEmail: email,
			expectID:    userID,
		},
		{
			name:        "expired token",
			token:       expiredTokenString,
			expectError: true,
		},
		{
			name:        "wrong signature",
			token:       wrongSignatureToken,
			expectError: true,
		},
		{
			name:        "invalid token format",
			token:       "invalid.token.format",
			expectError: true,
		},
		{
			name:        "empty token",
			token:       "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := service.ValidateToken(tt.token)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if claims == nil {
				t.Error("expected claims but got nil")
				return
			}

			if claims.Email != tt.expectEmail {
				t.Errorf("expected email %s, got %s", tt.expectEmail, claims.Email)
			}

			if claims.UserID != tt.expectID {
				t.Errorf("expected user ID %s, got %s", tt.expectID, claims.UserID)
			}
		})
	}
}

func TestService_ValidateToken_TokenExpiration(t *testing.T) {
	// Create a service with very short token duration for testing
	cfg, err := config.InitConfig("../../../config.test.yaml")
	secret := "test-secret"
	if err == nil {
		secret = cfg.Secret
	}
	shortDurationService := NewService(secret, 10*time.Millisecond, time.Hour)
	
	userID := uuid.New().String()
	email := "test@example.com"

	token, err := shortDurationService.GenerateToken(userID, email)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	// Wait for token to expire
	time.Sleep(20 * time.Millisecond)

	_, err = shortDurationService.ValidateToken(token)
	if err == nil {
		t.Error("expected error for expired token but got none")
	}
}

func TestService_TokenDuration(t *testing.T) {
	tokenDuration := 30 * time.Minute
	cfg, err := config.InitConfig("../../../config.test.yaml")
	secret := "test-secret"
	if err == nil {
		secret = cfg.Secret
	}
	service := NewService(secret, tokenDuration, time.Hour*24)
	
	userID := uuid.New().String()
	email := "test@example.com"

	token, err := service.GenerateToken(userID, email)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	claims, err := service.ValidateToken(token)
	if err != nil {
		t.Fatalf("failed to validate token: %v", err)
	}

	// Check if the token expires at the expected time (with some tolerance for execution time)
	expectedExpiry := time.Now().Add(tokenDuration)
	actualExpiry := claims.ExpiresAt.Time
	
	diff := actualExpiry.Sub(expectedExpiry)
	if diff < -time.Second || diff > time.Second {
		t.Errorf("expected expiry around %v, got %v (diff: %v)", expectedExpiry, actualExpiry, diff)
	}
}

func TestNewService(t *testing.T) {
	secretKey := "test-secret-key"
	tokenDuration := time.Hour
	refreshDuration := time.Hour * 24

	service := NewService(secretKey, tokenDuration, refreshDuration)

	if service == nil {
		t.Error("expected service instance but got nil")
	}

	if string(service.secretKey) != secretKey {
		t.Errorf("expected secret key %s, got %s", secretKey, string(service.secretKey))
	}

	if service.tokenDuration != tokenDuration {
		t.Errorf("expected token duration %v, got %v", tokenDuration, service.tokenDuration)
	}

	if service.refreshDuration != refreshDuration {
		t.Errorf("expected refresh duration %v, got %v", refreshDuration, service.refreshDuration)
	}
}

func TestClaims(t *testing.T) {
	userID := uuid.New().String()
	email := "test@example.com"
	jti := uuid.New().String()
	
	claims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jti,
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}

	if claims.UserID != userID {
		t.Errorf("expected UserID %s, got %s", userID, claims.UserID)
	}

	if claims.Email != email {
		t.Errorf("expected Email %s, got %s", email, claims.Email)
	}

	if claims.ID != jti {
		t.Errorf("expected ID %s, got %s", jti, claims.ID)
	}

	if claims.Subject != userID {
		t.Errorf("expected Subject %s, got %s", userID, claims.Subject)
	}
}
