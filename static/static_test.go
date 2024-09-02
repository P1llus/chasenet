package static

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetStyles(t *testing.T) {
	// TestGetStyles tests the GetStyles function
	// It checks if the GetStyles function returns the correct value
	// It checks if the returned value is not nil
	// It checks if the returned value is not empty
	// It checks if the returned value is not a string
	fs := GetStyles()
	output, err := fs.ReadFile("styles.css")
	if err != nil {
		t.Errorf("GetStyles() = %v; want not nil", output)
	}
	assert.Contains(t, string(output), `@import "https://cdn.jsdelivr.net/npm/bulma@1.0.2/css/bulma.min.css"`)
}
