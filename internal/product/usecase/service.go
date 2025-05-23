package usecase

import "mlsport/internal/product/domain"

type ProductService struct {
	Repo domain.ProductRepository
}

func NewProductService(repo domain.ProductRepository) *ProductService {
	return &ProductService{Repo: repo}
}

func (s *ProductService) Create(p *domain.Product) error {
	return s.Repo.Create(p)
}

func (s *ProductService) GetAll() ([]domain.Product, error) {
	return s.Repo.FindAll()
}

func (s *ProductService) GetByID(id string) (*domain.Product, error) {
	return s.Repo.FindByID(id)
}
func (s *ProductService) GetCategories() ([]string, error) {
	return s.Repo.GetCategories()
}
func (s *ProductService) GetByCategory(cat string) ([]domain.Product, error) {
	return s.Repo.FindByCategory(cat)
}

func (s *ProductService) Update(p *domain.Product) error {
	return s.Repo.Update(p)
}

func (s *ProductService) Patch(id string, fields map[string]interface{}) error {
	return s.Repo.Patch(id, fields)
}

func (s *ProductService) Delete(id string) error {
	return s.Repo.Delete(id)
}
func (s *ProductService) GetMetrics() (map[string]interface{}, error) {
	return s.Repo.GetMetrics()
}
