package delivery

import (
	"mlsport/internal/product/domain"
	"mlsport/internal/product/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	Service *usecase.ProductService
}

func NewProductHandler(s *usecase.ProductService) *ProductHandler {
	return &ProductHandler{Service: s}
}

// GetAll godoc
// @Summary Obtener todos los productos
// @Description Devuelve una lista completa de productos registrados en el sistema.
// @Tags Productos
// @Produce json
// @Success 200 {array} domain.Product
// @Router /products [get]
func (h *ProductHandler) GetAll(c *gin.Context) {
	products, err := h.Service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error obteniendo productos"})
		return
	}

	if len(products) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "no hay productos disponibles",
			"data":    []interface{}{},
		})
		return
	}

	c.JSON(http.StatusOK, products)
}

// GetByID godoc
// @Summary Consultar un producto por ID
// @Description Retorna la información detallada de un producto específico.
// @Tags Productos
// @Produce json
// @Param id path string true "ID del producto"
// @Success 200 {object} domain.Product
// @Failure 404 {object} map[string]string
// @Router /products/{id} [get]
func (h *ProductHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	product, err := h.Service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "producto no encontrado"})
		return
	}
	c.JSON(http.StatusOK, product)
}

// Create godoc
// @Summary Crear nuevo producto
// @Description Permite registrar un nuevo producto en la base de datos.
// @Tags Productos
// @Accept json
// @Produce json
// @Param producto body domain.Product true "Producto a registrar"
// @Success 201 {object} domain.Product
// @Failure 400 {object} map[string]string
// @Router /products [post]
func (h *ProductHandler) Create(c *gin.Context) {
	var input domain.Product
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "formato inválido"})
		return
	}
	err := h.Service.Create(&input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo crear el producto"})
		return
	}
	c.JSON(http.StatusCreated, input)
}

// Update godoc
// @Summary Reemplazar producto existente
// @Description Actualiza todos los campos de un producto existente con los nuevos valores proporcionados.
// @Tags Productos
// @Accept json
// @Produce json
// @Param id path string true "ID del producto"
// @Param producto body domain.Product true "Datos actualizados del producto"
// @Success 200 {object} domain.Product
// @Failure 400 {object} map[string]string
// @Router /products/{id} [put]
func (h *ProductHandler) Update(c *gin.Context) {
	var input domain.Product
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "formato inválido"})
		return
	}
	input.ID = c.Param("id")
	err := h.Service.Update(&input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo actualizar"})
		return
	}
	c.JSON(http.StatusOK, input)
}

// Patch godoc
// @Summary Actualizar parcialmente un producto
// @Description Permite actualizar solo algunos campos de un producto existente.
// @Tags Productos
// @Accept json
// @Produce json
// @Param id path string true "ID del producto"
// @Param fields body object true "Campos a modificar (por ejemplo: stock, price)"
// @Success 200 {object} map[string]string
// @Router /products/{id} [patch]
func (h *ProductHandler) Patch(c *gin.Context) {
	id := c.Param("id")
	var fields map[string]interface{}
	if err := c.ShouldBindJSON(&fields); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "formato inválido"})
		return
	}
	err := h.Service.Patch(id, fields)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo aplicar el patch"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "actualizado"})
}

// Delete godoc
// @Summary Eliminar producto
// @Description Borra un producto según su ID.
// @Tags Productos
// @Param id path string true "ID del producto"
// @Success 200 {object} map[string]string
// @Router /products/{id} [delete]
func (h *ProductHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	err := h.Service.Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo eliminar"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "eliminado"})
}

// GetByCategory godoc
// @Summary Obtener productos por categoría
// @Description Filtra los productos según la categoría proporcionada.
// @Tags Productos
// @Produce json
// @Param category path string true "Nombre de la categoría"
// @Success 200 {array} domain.Product
// @Router /products/categories/{category} [get]
func (h *ProductHandler) GetByCategory(c *gin.Context) {
	cat := c.Param("category")
	products, err := h.Service.GetByCategory(cat)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudieron filtrar los productos"})
		return
	}
	c.JSON(http.StatusOK, products)
}

// GetCategories godoc
// @Summary Listar categorías únicas
// @Description Retorna una lista de categorías derivadas de los productos registrados.
// @Tags Productos
// @Produce json
// @Success 200 {array} string
// @Router /products/categories [get]
func (h *ProductHandler) GetCategories(c *gin.Context) {
	list, err := h.Service.GetCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudieron obtener las categorías"})
		return
	}

	if len(list) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "no hay categorías", "data": []string{}})
		return
	}

	c.JSON(http.StatusOK, list)
}

// GetMetrics godoc
// @Summary Métricas de productos
// @Description Devuelve métricas agregadas como total de productos, promedio de precios y stock acumulado.
// @Tags Productos
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /products/metrics [get]
func (h *ProductHandler) GetMetrics(c *gin.Context) {
	data, err := h.Service.GetMetrics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudieron calcular las métricas"})
		return
	}
	c.JSON(http.StatusOK, data)
}

// GetDashboard godoc
// @Summary Dashboard de productos y métricas
// @Description Consulta combinada que obtiene productos y métricas en paralelo.
// @Tags Productos
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /products/dashboard [get]
func (h *ProductHandler) GetDashboard(c *gin.Context) {
	type result struct {
		Products []domain.Product
		Metrics  map[string]interface{}
		Err      error
	}

	productCh := make(chan []domain.Product)
	metricCh := make(chan map[string]interface{})
	errorCh := make(chan error, 2)

	go func() {
		list, err := h.Service.GetAll()
		if err != nil {
			errorCh <- err
			return
		}
		productCh <- list
	}()

	go func() {
		metrics, err := h.Service.GetMetrics()
		if err != nil {
			errorCh <- err
			return
		}
		metricCh <- metrics
	}()

	var res result

	for i := 0; i < 2; i++ {
		select {
		case products := <-productCh:
			res.Products = products
		case metrics := <-metricCh:
			res.Metrics = metrics
		case err := <-errorCh:
			res.Err = err
		}
	}

	if res.Err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": res.Err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
