package controllers

import (
	"net/http"
	"strconv"

	"go_app/models"

	"github.com/labstack/echo/v4"
)

var carts []models.Cart
var cartID uint = 1

// Utworzenie koszyka
func CreateCart(c echo.Context) error {
    cart := new(models.Cart)
    if err := c.Bind(cart); err != nil {
        return c.JSON(http.StatusBadRequest, err)
    }
    cart.ID = cartID
    cartID++
    carts = append(carts, *cart)
    return c.JSON(http.StatusCreated, cart)
}

// Pobranie koszyka po ID
func GetCart(c echo.Context) error {
    id, _ := strconv.Atoi(c.Param("id"))
    for _, cart := range carts {
        if cart.ID == uint(id) {
            return c.JSON(http.StatusOK, cart)
        }
    }
    return c.JSON(http.StatusNotFound, "Cart not found")
}
