package service_test

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"productservice/pkg/common/domain"
	"productservice/pkg/product/domain/model"
	"productservice/pkg/product/domain/service"
)

// --- Mocks ---
type mockProductRepository struct {
	products map[uuid.UUID]model.Product
	// Позволяет нам "подсунуть" ошибку в тест
	simulateError error
}

func newMockProductRepository() *mockProductRepository {
	return &mockProductRepository{
		products: make(map[uuid.UUID]model.Product),
	}
}

func (m *mockProductRepository) NextID() (uuid.UUID, error) {
	if m.simulateError != nil {
		return uuid.Nil, m.simulateError
	}
	return uuid.New(), nil
}

func (m *mockProductRepository) Store(product model.Product) error {
	if m.simulateError != nil {
		return m.simulateError
	}
	m.products[product.ProductID] = product
	return nil
}

func (m *mockProductRepository) Find(spec model.FindSpec) (*model.Product, error) {
	if m.simulateError != nil {
		return nil, m.simulateError
	}
	if spec.ProductID != nil {
		if p, ok := m.products[*spec.ProductID]; ok {
			return &p, nil
		}
	}
	if spec.Name != nil {
		for _, p := range m.products {
			if p.Name == *spec.Name {
				return &p, nil
			}
		}
	}
	return nil, model.ErrProductNotFound
}

func (m *mockProductRepository) HardDelete(productID uuid.UUID) error {
	if m.simulateError != nil {
		return m.simulateError
	}
	delete(m.products, productID)
	return nil
}

type mockEventDispatcher struct {
	dispatchedEvents []domain.Event
}

func (m *mockEventDispatcher) Dispatch(event domain.Event) error {
	m.dispatchedEvents = append(m.dispatchedEvents, event)
	return nil
}

// --- Tests ---

func TestProductService_CreateProduct_Success(t *testing.T) {
	// Arrange
	repo := newMockProductRepository()
	dispatcher := &mockEventDispatcher{}
	productService := service.NewProductService(repo, dispatcher)
	productName := "Test Coffee"
	productPrice := int64(300)

	// Act
	productID, err := productService.CreateProduct(productName, productPrice)

	// Assert
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, productID)

	// Проверяем, что продукт сохранился в нашем фейковом репозитории
	storedProduct, ok := repo.products[productID]
	require.True(t, ok)
	assert.Equal(t, productName, storedProduct.Name)
	assert.Equal(t, productPrice, storedProduct.Price)

	// Проверяем, что было отправлено правильное событие
	require.Len(t, dispatcher.dispatchedEvents, 1)
	createdEvent, ok := dispatcher.dispatchedEvents[0].(*model.ProductCreated)
	require.True(t, ok)
	assert.Equal(t, productID, createdEvent.ProductID)
	assert.Equal(t, productName, createdEvent.Name)
	assert.WithinDuration(t, time.Now(), createdEvent.CreatedAt, time.Second)
}

func TestProductService_CreateProduct_NameAlreadyExists(t *testing.T) {
	// Arrange
	repo := newMockProductRepository()
	dispatcher := &mockEventDispatcher{}
	productService := service.NewProductService(repo, dispatcher)

	// Заранее "кладем" продукт в репозиторий
	existingProductName := "Existing Tea"
	repo.products[uuid.New()] = model.Product{Name: existingProductName}

	// Act
	_, err := productService.CreateProduct(existingProductName, 500)

	// Assert
	require.Error(t, err)
	assert.True(t, errors.Is(err, model.ErrProductNameAlreadyUsed))
	assert.Empty(t, dispatcher.dispatchedEvents)
}
