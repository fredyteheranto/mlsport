// @title mlsport API
// @version 1.0
// @description API REST para productos deportivos.
// @host localhost:8080
// @BasePath /api

package main

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"mlsport/config"
	_ "mlsport/docs"
	"mlsport/internal/product/delivery"
	"mlsport/internal/product/infrastructure"
	"mlsport/internal/product/usecase"
)

func main() {
	config.InitMongo()

	repo := infrastructure.NewMongoProductRepo()
	service := usecase.NewProductService(repo)
	handler := delivery.NewProductHandler(service)

	r := gin.Default()

	// Swagger UI en /swagger/index.html
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Ruta raíz
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Bienvenido a la API de mlsport 🎽",
			"docs":    "Visita /swagger para ver la documentación completa",
		})
	})

	// Ruta base de API
	api := r.Group("/api")
	{
		api.GET("", func(c *gin.Context) {
			c.Redirect(302, "/swagger/index.html")
		})

		api.GET("/products/dashboard", handler.GetDashboard)

		products := api.Group("/products")
		{
			products.GET("", handler.GetAll)
			products.GET("/:id", handler.GetByID)
			products.GET("/categories/:category", handler.GetByCategory)
			products.GET("/metrics", handler.GetMetrics)
			products.GET("/categories", handler.GetCategories)

			products.POST("", handler.Create)
			products.PUT("/:id", handler.Update)
			products.PATCH("/:id", handler.Patch)
			products.DELETE("/:id", handler.Delete)
		}
	}

	if err := r.SetTrustedProxies(nil); err != nil {
		log.Fatalf("Error setting trusted proxies: %v", err)
	}

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Error running server: %v", err)
	}

}
