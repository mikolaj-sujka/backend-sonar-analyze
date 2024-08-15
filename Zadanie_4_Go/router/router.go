package router

import (
	"go_app/controllers"

	"github.com/labstack/echo/v4"
)

func New() *echo.Echo {
	return echo.New()
}

func Configure(e *echo.Echo)  {
    // Trasy dla produktów
    e.POST("/products", controllers.CreateProduct)
    e.GET("/products", controllers.GetProducts)
    e.GET("/products/:id", controllers.GetProduct)
    e.PUT("/products/:id", controllers.UpdateProduct)
    e.DELETE("/products/:id", controllers.DeleteProduct)

    // Trasy dla koszyków
    e.POST("/carts", controllers.CreateCart)
    e.GET("/carts/:id", controllers.GetCart)

    // Trasy dla kategorii
    e.POST("/categories", controllers.CreateCategory)
    e.GET("/categories", controllers.GetCategories)
}
