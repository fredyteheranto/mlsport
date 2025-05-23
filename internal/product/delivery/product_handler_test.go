package delivery

import (
	"bytes"
	"encoding/json"
	"errors"
	"mlsport/internal/product/domain"
	"mlsport/internal/product/usecase"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockRepo struct{}

func (m *mockRepo) Create(p *domain.Product) error { return nil }
func (m *mockRepo) FindAll() ([]domain.Product, error) {
	return []domain.Product{{ID: "1", Name: "Balón"}}, nil
}
func (m *mockRepo) FindByID(id string) (*domain.Product, error) {
	return &domain.Product{ID: id, Name: "Balón"}, nil
}
func (m *mockRepo) FindByCategory(cat string) ([]domain.Product, error) {
	return []domain.Product{{Category: cat}}, nil
}
func (m *mockRepo) Update(p *domain.Product) error                       { return nil }
func (m *mockRepo) Patch(id string, fields map[string]interface{}) error { return nil }
func (m *mockRepo) Delete(id string) error                               { return nil }
func (m *mockRepo) GetMetrics() (map[string]interface{}, error) {
	return map[string]interface{}{"total_products": 1}, nil
}
func (m *mockRepo) GetCategories() ([]string, error) { return []string{"Ropa", "Calzado"}, nil }

func newMockHandler() *ProductHandler {
	repo := &mockRepo{}
	service := usecase.NewProductService(repo)
	return NewProductHandler(service)
}

func TestGetAllHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := newMockHandler()

	req, _ := http.NewRequest("GET", "/api/products", nil)
	resp := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(resp)
	c.Request = req

	handler.GetAll(c)

	assert.Equal(t, http.StatusOK, resp.Code)

	var products []domain.Product
	err := json.Unmarshal(resp.Body.Bytes(), &products)
	assert.NoError(t, err)
	assert.Len(t, products, 1)
	assert.Equal(t, "Balón", products[0].Name)
}

func TestGetByIDHandler_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := newMockHandler()

	req, _ := http.NewRequest("GET", "/api/products/999", nil)
	resp := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(resp)
	c.Params = []gin.Param{{Key: "id", Value: "999"}}
	c.Request = req

	// Modificamos mock para simular no encontrado
	handler.Service = usecase.NewProductService(&notFoundMockRepo{})

	handler.GetByID(c)

	assert.Equal(t, http.StatusNotFound, resp.Code)
}

type notFoundMockRepo struct{}

func (m *notFoundMockRepo) Create(p *domain.Product) error     { return nil }
func (m *notFoundMockRepo) FindAll() ([]domain.Product, error) { return nil, nil }
func (m *notFoundMockRepo) FindByID(id string) (*domain.Product, error) {
	return nil, errors.New("producto no encontrado")
}
func (m *notFoundMockRepo) FindByCategory(cat string) ([]domain.Product, error)  { return nil, nil }
func (m *notFoundMockRepo) Update(p *domain.Product) error                       { return nil }
func (m *notFoundMockRepo) Patch(id string, fields map[string]interface{}) error { return nil }
func (m *notFoundMockRepo) Delete(id string) error                               { return nil }
func (m *notFoundMockRepo) GetMetrics() (map[string]interface{}, error)          { return nil, nil }
func (m *notFoundMockRepo) GetCategories() ([]string, error)                     { return nil, nil }

func TestCreateHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := newMockHandler()

	product := domain.Product{Name: "Nuevo"}
	jsonValue, _ := json.Marshal(product)
	req, _ := http.NewRequest("POST", "/api/products", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(resp)
	c.Request = req

	handler.Create(c)

	assert.Equal(t, http.StatusCreated, resp.Code)
}
