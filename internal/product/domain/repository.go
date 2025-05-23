package domain

type ProductRepository interface {
	Create(product *Product) error
	FindAll() ([]Product, error)
	FindByID(id string) (*Product, error)
	FindByCategory(category string) ([]Product, error)
	Update(product *Product) error
	Patch(id string, fields map[string]interface{}) error
	Delete(id string) error
	GetMetrics() (map[string]interface{}, error)
	GetCategories() ([]string, error)
}
