package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrProductNotFound        = errors.New("product not found")
	ErrProductNameAlreadyUsed = errors.New("product name already used")
)

// Product представляет доменную модель продукта
type Product struct {
	ProductID uuid.UUID
	Name      string
	Price     int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

type FindSpec struct {
	ProductID *uuid.UUID
	Name      *string
}

type ProductRepository interface {
	NextID() (uuid.UUID, error)
	Store(product Product) error
	Find(spec FindSpec) (*Product, error)
	HardDelete(productID uuid.UUID) error
}
