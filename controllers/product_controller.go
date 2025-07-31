package controllers

import (
	"net/http"
	"test1/config"
	"test1/models"

	"github.com/labstack/echo/v4"
)

func CreateProduct(c echo.Context) error {
	var product models.Product
	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}
	if err := config.DB.Create(&product).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to create product"})
	}
	return c.JSON(http.StatusCreated, product)
}

func GetAllProducts(c echo.Context) error {
	var products []models.Product
	if err := config.DB.Find(&products).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to fetch products"})
	}
	return c.JSON(http.StatusOK, products)
}

func GetProductByID(c echo.Context) error {
	id := c.Param("id")
	var product models.Product
	if err := config.DB.First(&product, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Product not found"})
	}
	return c.JSON(http.StatusOK, product)
}

func UpdateProduct(c echo.Context) error {
	id := c.Param("id")
	var existing models.Product
	if err := config.DB.First(&existing, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Product not found"})
	}
	var updateData models.Product
	if err := c.Bind(&updateData); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}
	updateData.ID = existing.ID
	if err := config.DB.Save(&updateData).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to update product"})
	}
	return c.JSON(http.StatusOK, updateData)
}

func DeleteProduct(c echo.Context) error {
	id := c.Param("id")
	var product models.Product
	if err := config.DB.First(&product, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Product not found"})
	}
	if err := config.DB.Delete(&product).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to delete product"})
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "Product deleted successfully"})
}
