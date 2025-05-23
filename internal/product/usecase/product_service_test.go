package usecase

import (
	"errors"
	"testing"

	"mlsport/internal/product/domain"

	"github.com/stretchr/testify/assert"
)

type mockRepo struct{}

func (m *mockRepo) Create(p *domain.Product) error { return nil }
func (m *mockRepo) FindAll() ([]domain.Product, error) {
	return []domain.Product{{Name: "Balón"}}, nil
}
func (m *mockRepo) FindByID(id string) (*domain.Product, error) {
	return &domain.Product{ID: id, Name: "Zapatilla"}, nil
}
func (m *mockRepo) FindByCategory(cat string) ([]domain.Product, error) {
	return []domain.Product{{Category: cat}}, nil
}
func (m *mockRepo) Update(p *domain.Product) error                       { return nil }
func (m *mockRepo) Patch(id string, fields map[string]interface{}) error { return nil }
func (m *mockRepo) Delete(id string) error                               { return nil }
func (m *mockRepo) GetMetrics() (map[string]interface{}, error) {
	return map[string]interface{}{
		"total_products": 3,
		"top_categories": []string{"Ropa", "Calzado"},
		"total_stock":    200,
		"average_price":  89.5,
	}, nil
}
func (m *mockRepo) GetCategories() ([]string, error) { return []string{"Ropa", "Calzado"}, nil }

// mockRepo que simula errores
type errorMockRepo struct{}

func (m *errorMockRepo) Create(p *domain.Product) error { return errors.New("error simulado create") }
func (m *errorMockRepo) FindAll() ([]domain.Product, error) {
	return nil, errors.New("error simulado findall")
}
func (m *errorMockRepo) FindByID(id string) (*domain.Product, error) {
	return nil, errors.New("error simulado findbyid")
}
func (m *errorMockRepo) FindByCategory(cat string) ([]domain.Product, error) {
	return nil, errors.New("error simulado findbycategory")
}
func (m *errorMockRepo) Update(p *domain.Product) error { return errors.New("error simulado update") }
func (m *errorMockRepo) Patch(id string, fields map[string]interface{}) error {
	return errors.New("error simulado patch")
}
func (m *errorMockRepo) Delete(id string) error { return errors.New("error simulado delete") }
func (m *errorMockRepo) GetMetrics() (map[string]interface{}, error) {
	return nil, errors.New("error simulado metrics")
}
func (m *errorMockRepo) GetCategories() ([]string, error) {
	return nil, errors.New("error simulado categories")
}

func TestCreateProduct(t *testing.T) {
	repo := &mockRepo{}
	service := NewProductService(repo)

	err := service.Create(&domain.Product{Name: "Nuevo Producto"})

	assert.NoError(t, err)
}

func TestGetAll(t *testing.T) {
	repo := &mockRepo{}
	service := NewProductService(repo)

	list, err := service.GetAll()

	assert.NoError(t, err)
	assert.Len(t, list, 1)
	assert.Equal(t, "Balón", list[0].Name)
}

func TestGetByID(t *testing.T) {
	repo := &mockRepo{}
	service := NewProductService(repo)

	product, err := service.GetByID("123")

	assert.NoError(t, err)
	assert.Equal(t, "123", product.ID)
	assert.Equal(t, "Zapatilla", product.Name)
}

func TestGetByCategory(t *testing.T) {
	repo := &mockRepo{}
	service := NewProductService(repo)

	prods, err := service.GetByCategory("Accesorios")

	assert.NoError(t, err)
	assert.Equal(t, "Accesorios", prods[0].Category)
}

func TestGetMetrics(t *testing.T) {
	repo := &mockRepo{}
	service := NewProductService(repo)

	data, err := service.GetMetrics()

	assert.NoError(t, err)
	assert.Equal(t, 3, data["total_products"])
	assert.Contains(t, data["top_categories"], "Ropa")
	assert.Equal(t, 200, data["total_stock"])
	assert.Equal(t, 89.5, data["average_price"])
}

func TestUpdateProduct(t *testing.T) {
	repo := &mockRepo{}
	service := NewProductService(repo)

	p := &domain.Product{ID: "123", Name: "Nuevo nombre"}
	err := service.Update(p)

	assert.NoError(t, err)
}

func TestUpdateProduct_Error(t *testing.T) {
	repo := &errorMockRepo{}
	service := NewProductService(repo)

	p := &domain.Product{ID: "123", Name: "Error"}
	err := service.Update(p)

	assert.Error(t, err)
	assert.Equal(t, "error simulado update", err.Error())
}

func TestPatchProduct(t *testing.T) {
	repo := &mockRepo{}
	service := NewProductService(repo)

	err := service.Patch("123", map[string]interface{}{"price": 99.9})
	assert.NoError(t, err)
}

func TestPatchProduct_Error(t *testing.T) {
	repo := &errorMockRepo{}
	service := NewProductService(repo)

	err := service.Patch("123", map[string]interface{}{"price": 99.9})
	assert.Error(t, err)
	assert.Equal(t, "error simulado patch", err.Error())
}

func TestDeleteProduct(t *testing.T) {
	repo := &mockRepo{}
	service := NewProductService(repo)

	err := service.Delete("123")
	assert.NoError(t, err)
}

func TestDeleteProduct_Error(t *testing.T) {
	repo := &errorMockRepo{}
	service := NewProductService(repo)

	err := service.Delete("123")
	assert.Error(t, err)
	assert.Equal(t, "error simulado delete", err.Error())
}
