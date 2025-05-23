package delivery

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"mlsport/internal/product/domain"
	"mlsport/internal/product/usecase"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockDashboardRepo struct{}

func (m *mockDashboardRepo) FindAll() ([]domain.Product, error) {
	return []domain.Product{
		{ID: "1", Name: "Balón"},
		{ID: "2", Name: "Guayos"},
	}, nil
}
func (m *mockDashboardRepo) GetMetrics() (map[string]interface{}, error) {
	return map[string]interface{}{
		"total_products": 2,
		"total_stock":    100,
		"average_price":  50.5,
		"top_categories": []string{"Ropa", "Calzado"},
	}, nil
}

// Implementa otros métodos vacíos para cumplir la interfaz:
func (m *mockDashboardRepo) Create(p *domain.Product) error                       { return nil }
func (m *mockDashboardRepo) FindByID(id string) (*domain.Product, error)          { return nil, nil }
func (m *mockDashboardRepo) FindByCategory(cat string) ([]domain.Product, error)  { return nil, nil }
func (m *mockDashboardRepo) Update(p *domain.Product) error                       { return nil }
func (m *mockDashboardRepo) Patch(id string, fields map[string]interface{}) error { return nil }
func (m *mockDashboardRepo) Delete(id string) error                               { return nil }
func (m *mockDashboardRepo) GetCategories() ([]string, error)                     { return nil, nil }

func newMockDashboardHandler() *ProductHandler {
	repo := &mockDashboardRepo{}
	svc := usecase.NewProductService(repo)
	return NewProductHandler(svc)
}

func TestGetDashboardSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := newMockDashboardHandler()

	req, _ := http.NewRequest("GET", "/api/dashboard", nil)
	resp := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(resp)
	c.Request = req

	handler.GetDashboard(c)

	assert.Equal(t, http.StatusOK, resp.Code)
}
