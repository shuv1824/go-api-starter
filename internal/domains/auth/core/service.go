package core

import "github.com/google/uuid"

type Service interface {
	GenerateToken(userID uuid.UUID, email string) (string, error)
	ValidateToken(token string) (*Claims, error)
}
