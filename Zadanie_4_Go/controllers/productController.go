package controllers

import (
	"go_app/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

var products []models.Product
var productID uint = 1

// Utworzenie produktu
func CreateProduct(c echo.Context) error {
    product := new(models.Product)
    if err := c.Bind(product); err != nil {
        return c.JSON(http.StatusBadRequest, err)
    }
    product.ID = productID
    productID++
    products = append(products, *product)
    return c.JSON(http.StatusCreated, product)
}

// Pobranie wszystkich produktów
func GetProducts(c echo.Context) error {
    return c.JSON(http.StatusOK, products)
}

// Pobranie produktu po ID
func GetProduct(c echo.Context) error {
    id, _ := strconv.Atoi(c.Param("id"))
    for _, product := range products {
        if product.ID == uint(id) {
            return c.JSON(http.StatusOK, product)
        }
    }
    return c.JSON(http.StatusNotFound, "Product not found")
}

// Aktualizacja produktu
func UpdateProduct(c echo.Context) error {
    id, _ := strconv.Atoi(c.Param("id"))
    updatedProduct := new(models.Product)
    if err := c.Bind(updatedProduct); err != nil {
        return c.JSON(http.StatusBadRequest, err)
    }
    for i, product := range products {
        if product.ID == uint(id) {
            products[i] = *updatedProduct
            products[i].ID = uint(id)
            return c.JSON(http.StatusOK, products[i])
        }
    }
    return c.JSON(http.StatusNotFound, "Product not found")
}

// Usunięcie produktu
func DeleteProduct(c echo.Context) error {
    id, _ := strconv.Atoi(c.Param("id"))
    for i, product := range products {
        if product.ID == uint(id) {
            products = append(products[:i], products[i+1:]...)
            return c.NoContent(http.StatusNoContent)
        }
    }
    return c.JSON(http.StatusNotFound, "Product not found")
}
