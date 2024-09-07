package server

import (
	"github.com/P1llus/chasenet/pages/about"
	"github.com/P1llus/chasenet/pages/blog"
	"github.com/P1llus/chasenet/static"
	"github.com/labstack/echo/v4"
)

// Init all routes
func initRoutes(e *echo.Echo) {
	initAboutRoutes(e)
	initBlogRoutes(e)
	initStaticRoutes(e)
}

// Static routes
func initStaticRoutes(e *echo.Echo) {
	staticFS := static.GetStyles()
	e.StaticFS("/static", staticFS)
}

// About routes
func initAboutRoutes(e *echo.Echo) {
	aboutManager := about.NewAboutManager()
	err := aboutManager.LoadAboutPage()
	if err != nil {
		e.Logger.Fatal(err)
		return
	}
	e.GET("/about", getAboutHandler(&aboutManager))
}

// Blog routes
func initBlogRoutes(e *echo.Echo) {
	blogManager := blog.NewBlogManager()
	err := blogManager.LoadBlogPosts()
	if err != nil {
		e.Logger.Fatal(err)
		return
	}
	e.GET("/blog/:name", getBlogPostHandler(&blogManager))
	e.GET("/blog", getBlogPostsHandler(&blogManager))
	e.GET("/", getBlogPostsHandler(&blogManager))
}
