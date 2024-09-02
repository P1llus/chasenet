package server

import (
	"net/http"

	"github.com/P1llus/chasenet/pages/about"
	"github.com/P1llus/chasenet/pages/blog"
	"github.com/P1llus/chasenet/static"
	"github.com/P1llus/chasenet/views"
	"github.com/labstack/echo/v4"
)

func Start() {
	e := echo.New()
	staticFS := static.GetStyles()
	blogManager := blog.NewBlogManager()
	err := blogManager.LoadBlogPosts()
	if err != nil {
		e.Logger.Fatal(err)
		return
	}
	aboutManager := about.NewAboutManager()
	err = aboutManager.LoadAboutPage()
	if err != nil {
		e.Logger.Fatal(err)
		return
	}
	e.GET("/blog/:name", getPost(blogManager))
	e.GET("/blog", getPosts(blogManager))
	e.GET("/about", getAbout(aboutManager))
	e.StaticFS("/static", staticFS)
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":8080"))
}

func getAbout(pm about.AboutManager) echo.HandlerFunc {
	return func(c echo.Context) error {
		page := pm.GetAboutPage()
		return views.AboutPage(page).Render(
			c.Request().Context(),
			c.Response().Writer,
		)
	}
}

func getPost(pm blog.BlogManager) echo.HandlerFunc {
	return func(c echo.Context) error {
		post := pm.GetBlogPostBySlug(c.Param("name"))
		if post == nil {
			return c.String(http.StatusNotFound, "No blogpost with that name was found")
		}
		return views.PostPage(post).Render(
			c.Request().Context(),
			c.Response().Writer,
		)
	}
}

func getPosts(pm blog.BlogManager) echo.HandlerFunc {
	return func(c echo.Context) error {
		posts := pm.ListBlogPosts()
		if posts == nil {
			return c.String(http.StatusNotFound, "No blogposts found")
		}
		return views.PostsPage(posts).Render(
			c.Request().Context(),
			c.Response().Writer,
		)
	}
}
