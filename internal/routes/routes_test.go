package routes

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestHealthCheck(t *testing.T) {
	// Switch Gin to test mode
	gin.SetMode(gin.TestMode)

	// Setup router & register routes
	r := gin.Default()
	r.GET("/healthcheck", HealthCheckHandler)

	// Create the mock request
	req, err := http.NewRequest(http.MethodGet, "/healthcheck", nil)
	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}

	// Create a response recorder
	record := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(record, req)
	fmt.Println(record.Body)

	// Check response code
	expectedCode := http.StatusOK
	if record.Code != expectedCode {
		t.Fatalf("Expected to get status %d but instead got %d\n", expectedCode, record.Code)
	}

	// Check the response body.
	expectedResponse := `OK`
	if record.Body.String() != expectedResponse {
		t.Errorf("Expected to get %s but instead got %s\n", expectedResponse, record.Body.String())
	}
}
