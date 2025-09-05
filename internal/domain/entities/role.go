package entities

import (
	"time"

	"github.com/google/uuid"
)

type Role struct {
	ID         uuid.UUID              `json:"id" db:"id"`
	DomainID   uuid.UUID              `json:"domain_id" db:"domain_id"`
	RoleName   string                 `json:"role_name" db:"role_name"`
	RoleClaims map[string]interface{} `json:"role_claims" db:"role_claims"`
	CreatedAt  time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time              `json:"updated_at" db:"updated_at"`
}
