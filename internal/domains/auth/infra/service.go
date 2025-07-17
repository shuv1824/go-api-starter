package infra

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/shuv1824/go-api-starter/internal/domains/auth/core"
)

type service struct {
	secretKey       []byte
	tokenDuration   time.Duration
	refreshDuration time.Duration
}

func New(secretKey string, tokenDuration, refreshDuration time.Duration) *service {
	return &service{
		secretKey:       []byte(secretKey),
		tokenDuration:   tokenDuration,
		refreshDuration: refreshDuration,
	}
}

func (s *service) GenerateToken(userID uuid.UUID, email string) (string, error) {
	claims := core.Claims{
		UserID: userID.String(),
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			Subject:   userID.String(),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.tokenDuration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

func (s *service) ValidateToken(tokenString string) (*core.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &core.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return s.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*core.Claims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}
