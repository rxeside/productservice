package model

import "github.com/google/uuid"

// Product DTO
type Product struct {
	ProductID uuid.UUID
	Name      string
	Price     int64
}
