package server

import (
	"net/http"

	"github.com/P1llus/chasenet/pages/about"
	"github.com/P1llus/chasenet/pages/blog"
	"github.com/P1llus/chasenet/views"
	"github.com/labstack/echo/v4"
)

// About handlers
func getAboutHandler(m *about.AboutManager) echo.HandlerFunc {
	return func(c echo.Context) error {
		page := m.GetAboutPage()
		return views.AboutPage(page).Render(
			c.Request().Context(),
			c.Response().Writer,
		)
	}
}

// Blog handlers
func getBlogPostBySlugHandler(m *blog.BlogManager) echo.HandlerFunc {
	return func(c echo.Context) error {
		post := m.GetBlogPostBySlug(c.Param("name"))
		if post == nil {
			return c.String(http.StatusNotFound, "No blogpost with that name was found")
		}
		return views.PostPage(post).Render(
			c.Request().Context(),
			c.Response().Writer,
		)
	}
}

func getBlogPostsByTagHandler(m *blog.BlogManager) echo.HandlerFunc {
	return func(c echo.Context) error {
		tag := c.Param("tag")
		posts := m.GetBlogPostsByTag(tag)
		if posts == nil {
			return c.String(http.StatusNotFound, "No blogposts with that tag was found")
		}
		return views.PostsByTagPage(posts, tag).Render(
			c.Request().Context(),
			c.Response().Writer,
		)
	}
}

func getBlogPostsHandler(m *blog.BlogManager) echo.HandlerFunc {
	return func(c echo.Context) error {
		posts := m.ListBlogPosts()
		if posts == nil {
			return c.String(http.StatusNotFound, "No blogposts found")
		}
		return views.PostsPage(posts).Render(
			c.Request().Context(),
			c.Response().Writer,
		)
	}
}
