package controllers

import (
	"go_app/controllers"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCreateProduct(t *testing.T) {
    e := echo.New()
    req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(`{"name":"Product1","price":100.0,"category_id":1}`))
    req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)

    if assert.NoError(t, controllers.CreateProduct(c)) {
        assert.Equal(t, http.StatusCreated, rec.Code)
        assert.Contains(t, rec.Body.String(), "Product1")
    }
}

func TestGetProducts(t *testing.T) {
    e := echo.New()
    req := httptest.NewRequest(http.MethodGet, "/products", nil)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)

    if assert.NoError(t, controllers.GetProducts(c)) {
        assert.Equal(t, http.StatusOK, rec.Code)
    }
}

