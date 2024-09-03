package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServer(t *testing.T) {
	e := initEcho()

	testCases := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Root path",
			method:         http.MethodGet,
			path:           "/",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "About Page",
			method:         http.MethodGet,
			path:           "/about",
			expectedStatus: http.StatusOK,
			expectedBody:   `Currently working at`,
		},
		{
			name:           "Blog List",
			method:         http.MethodGet,
			path:           "/blog",
			expectedStatus: http.StatusOK,
			expectedBody:   `List of all blog posts`,
		},
		{
			name:           "Blog Post",
			method:         http.MethodGet,
			path:           "/blog/test",
			expectedStatus: http.StatusOK,
			expectedBody:   `testpage`,
		},
		{
			name:           "Not found",
			method:         http.MethodGet,
			path:           "/nonexistent",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			if tc.expectedBody != "" {
				body, _ := io.ReadAll(rec.Body)
				assert.Contains(t, string(body), tc.expectedBody)
			}
		})
	}
}
