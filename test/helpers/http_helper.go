package helpers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// SetupGinTestMode sets Gin to test mode
func SetupGinTestMode() {
	gin.SetMode(gin.TestMode)
}

// CreateTestContext creates a test Gin context
func CreateTestContext() (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	return c, w
}

// MakeJSONRequest creates a test HTTP request with JSON body
func MakeJSONRequest(t *testing.T, method, url string, body interface{}) *http.Request {
	var bodyReader io.Reader
	if body != nil {
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("Failed to marshal request body: %v", err)
		}
		bodyReader = bytes.NewBuffer(jsonBytes)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req
}

// ParseJSONResponse parses JSON response from recorder
func ParseJSONResponse(t *testing.T, w *httptest.ResponseRecorder, v interface{}) {
	if err := json.Unmarshal(w.Body.Bytes(), v); err != nil {
		t.Fatalf("Failed to parse response: %v, body: %s", err, w.Body.String())
	}
}

// AssertStatusCode asserts HTTP status code
func AssertStatusCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected status code %d, got %d", expected, actual)
	}
}

// AssertJSONResponse asserts the response contains expected fields
func AssertJSONResponse(t *testing.T, body []byte, expectedFields map[string]interface{}) {
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	for key, expectedValue := range expectedFields {
		actualValue, exists := response[key]
		if !exists {
			t.Errorf("Expected field '%s' not found in response", key)
			continue
		}

		if actualValue != expectedValue {
			t.Errorf("Field '%s': expected %v, got %v", key, expectedValue, actualValue)
		}
	}
}
