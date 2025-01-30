package main

// Test generated using Keploy
import (
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
    "github.com/gin-gonic/gin"
)

func TestPutURL_InvalidPayload_Returns400(t *testing.T) {
    // Create a mock Gin context
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)

    // Test with an invalid payload
    invalidPayload := `{"invalid": "data"}`
    c.Request, _ = http.NewRequest("POST", "/url", strings.NewReader(invalidPayload))
    c.Request.Header.Set("Content-Type", "application/json")
    putURL(c)

    if w.Code != http.StatusBadRequest {
        t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
    }
}
