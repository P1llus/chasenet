package server

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Start() {
	e := initEcho()
	if os.Getenv("APP_ENV") == "production" {
		e.Pre(middleware.HTTPSNonWWWRedirect())
	}
	if os.Getenv("APP_ENV") != "production" {
		e.Pre(middleware.NonWWWRedirect())
		e.Use(middleware.Logger())
	}
	e.Logger.Fatal(e.Start(":8080"))
}

func initEcho() *echo.Echo {
	e := echo.New()
	initRoutes(e)
	return e
}
