package routes

import (
	"test1/controllers"

	"github.com/labstack/echo/v4"
)

func UserRoutes(e *echo.Echo) {
	e.POST("/signup", controllers.Signup)
	e.POST("/login", controllers.Login)
}
