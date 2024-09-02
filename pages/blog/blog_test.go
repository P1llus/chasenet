package blog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBlogManager(t *testing.T) {
	blogManager := NewBlogManager()
	blogManager.LoadBlogPosts()
	output := blogManager.ListBlogPosts()
	assert.Contains(t, output.Posts[0].Title, "testpage")

	output2 := blogManager.GetBlogPostBySlug("testing")
	assert.Contains(t, output2.Title, "testingpage")
}
