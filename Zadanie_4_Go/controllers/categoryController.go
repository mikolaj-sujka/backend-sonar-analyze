package controllers

import (
	"net/http"

	"go_app/models"

	"github.com/labstack/echo/v4"
)

var categories []models.Category
var categoryID uint = 1

// Utworzenie kategorii
func CreateCategory(c echo.Context) error {
    category := new(models.Category)
    if err := c.Bind(category); err != nil {
        return c.JSON(http.StatusBadRequest, err)
    }
    category.ID = categoryID
    categoryID++
    categories = append(categories, *category)
    return c.JSON(http.StatusCreated, category)
}

// Pobranie wszystkich kategorii
func GetCategories(c echo.Context) error {
    return c.JSON(http.StatusOK, categories)
}
