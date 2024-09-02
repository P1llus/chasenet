package server

import (
	"github.com/labstack/echo/v4"
)

func Start() {
	e := initEcho()
	e.Logger.Fatal(e.Start(":8080"))
}

func initEcho() *echo.Echo {
	e := echo.New()
	initRoutes(e)
	return e
}
