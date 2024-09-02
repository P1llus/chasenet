package about

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAboutManager(t *testing.T) {
	aboutManager := NewAboutManager()
	aboutManager.LoadAboutPage()
	output := aboutManager.GetAboutPage()
	assert.Contains(t, output.Title, "About Me")
	assert.Contains(t, output.Content, "Random facts")
}
