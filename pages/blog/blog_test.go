package blog

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBlogManager(t *testing.T) {
	blogManager := NewBlogManager()
	err := blogManager.LoadBlogPosts()
	assert.NoError(t, err)

	// Test GetBlogPostBySlug
	output := blogManager.GetBlogPostBySlug("testing")
	assert.Contains(t, output.Title, "testingpage")

	// Test ListBlogPosts and verify sorting
	posts := blogManager.ListBlogPosts()
	assert.NotEmpty(t, posts, "Blog posts should not be empty")

	// Check if the posts are sorted by date in descending order
	for i := 0; i < len(posts.Posts)-1; i++ {
		currentDate, err := time.Parse("02-Jan-2006", posts.Posts[i].Date)
		assert.NoError(t, err, "Error parsing date for post %d", i)

		nextDate, err := time.Parse("02-Jan-2006", posts.Posts[i+1].Date)
		assert.NoError(t, err, "Error parsing date for post %d", i+1)

		assert.True(t, currentDate.After(nextDate) || currentDate.Equal(nextDate),
			"Post at index %d should have a date after or equal to post at index %d", i, i+1)
	}
}
