package main

import (
	"test1/config"
	"test1/routes"

	"github.com/labstack/echo/v4"
)

func main() {
	config.ConnectDB()
	e := echo.New()
	routes.UserRoutes(e)
	routes.ProductRoutes(e)

	e.Logger.Fatal(e.Start(":8081"))
}
