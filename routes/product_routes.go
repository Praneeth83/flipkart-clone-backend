package routes

import (
	"test1/controllers"

	"github.com/labstack/echo/v4"
)

func ProductRoutes(e *echo.Echo) {
	e.POST("/products", controllers.CreateProduct)
	e.GET("/products", controllers.GetAllProducts)
	e.GET("/products/:id", controllers.GetProductByID)
	e.PUT("/products/:id", controllers.UpdateProduct)
	e.DELETE("/products/:id", controllers.DeleteProduct)
}
