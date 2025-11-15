package service

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"productservice/pkg/common/domain"
	"productservice/pkg/product/domain/model"
)

type ProductService interface {
	CreateProduct(name string, price int64) (uuid.UUID, error)
	UpdateProduct(productID uuid.UUID, name string, price int64) error
	DeleteProduct(productID uuid.UUID) error
}

func NewProductService(
	productRepository model.ProductRepository,
	eventDispatcher domain.EventDispatcher,
) ProductService {
	return &productService{
		productRepository: productRepository,
		eventDispatcher:   eventDispatcher,
	}
}

type productService struct {
	productRepository model.ProductRepository
	eventDispatcher   domain.EventDispatcher
}

func (s *productService) CreateProduct(name string, price int64) (uuid.UUID, error) {
	_, err := s.productRepository.Find(model.FindSpec{
		Name: &name,
	})
	if err != nil && !errors.Is(err, model.ErrProductNotFound) {
		return uuid.Nil, err
	}
	if err == nil {
		return uuid.Nil, model.ErrProductNameAlreadyUsed
	}

	productID, err := s.productRepository.NextID()
	if err != nil {
		return uuid.Nil, err
	}

	currentTime := time.Now()
	err = s.productRepository.Store(model.Product{
		ProductID: productID,
		Name:      name,
		Price:     price,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	})
	if err != nil {
		return uuid.Nil, err
	}

	return productID, s.eventDispatcher.Dispatch(&model.ProductCreated{
		ProductID: productID,
		Name:      name,
		Price:     price,
		CreatedAt: currentTime,
	})
}

func (s *productService) UpdateProduct(productID uuid.UUID, name string, price int64) error {
	product, err := s.productRepository.Find(model.FindSpec{
		ProductID: &productID,
	})
	if err != nil {
		return err
	}

	if product.Name == name && product.Price == price {
		return nil // Нет изменений
	}

	// Проверка, что новое имя не занято другим продуктом
	if product.Name != name {
		existing, err := s.productRepository.Find(model.FindSpec{Name: &name})
		if err != nil && !errors.Is(err, model.ErrProductNotFound) {
			return err
		}
		if existing != nil && existing.ProductID != productID {
			return model.ErrProductNameAlreadyUsed
		}
	}

	currentTime := time.Now()
	product.Name = name
	product.Price = price
	product.UpdatedAt = currentTime

	err = s.productRepository.Store(*product)
	if err != nil {
		return err
	}

	updatedEvent := &model.ProductUpdated{
		ProductID: productID,
		UpdatedAt: currentTime,
	}
	updatedEvent.UpdatedFields.Name = &name
	updatedEvent.UpdatedFields.Price = &price

	return s.eventDispatcher.Dispatch(updatedEvent)
}

func (s *productService) DeleteProduct(productID uuid.UUID) error {
	_, err := s.productRepository.Find(model.FindSpec{
		ProductID: &productID,
	})
	if err != nil {
		return err
	}

	err = s.productRepository.HardDelete(productID)
	if err != nil {
		return err
	}

	return s.eventDispatcher.Dispatch(&model.ProductDeleted{
		ProductID: productID,
		DeletedAt: time.Now(),
	})
}
