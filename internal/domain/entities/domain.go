package entities

import "github.com/google/uuid"

type Domain struct {
	DomainID uuid.UUID `json:"domain_id" db:"domain_id"`
	Name     string    `json:"name" db:"name"`
	Domain   string    `json:"domain" db:"domain"`
}
