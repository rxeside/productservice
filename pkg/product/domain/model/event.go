package model

import (
	"time"

	"github.com/google/uuid"
)

// ProductCreated событие о создании продукта
type ProductCreated struct {
	ProductID uuid.UUID
	Name      string
	Price     int64 // Цена в копейках
	CreatedAt time.Time
}

func (e ProductCreated) Type() string {
	return "product_created"
}

// ProductUpdated событие об обновлении продукта
type ProductUpdated struct {
	ProductID     uuid.UUID
	UpdatedFields struct {
		Name  *string
		Price *int64
	}
	UpdatedAt time.Time
}

func (e ProductUpdated) Type() string {
	return "product_updated"
}

// ProductDeleted событие об удалении продукта
type ProductDeleted struct {
	ProductID uuid.UUID
	DeletedAt time.Time
}

func (e ProductDeleted) Type() string {
	return "product_deleted"
}
