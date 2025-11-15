package model

import "github.com/google/uuid"

// Product это модель данных для application слоя (DTO)
type Product struct {
	ProductID uuid.UUID
	Name      string
	Price     int64
}
