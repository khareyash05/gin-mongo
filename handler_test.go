package main

import (
    "net/http"
    "net/http/httptest"
    "testing"
    "github.com/gin-gonic/gin"
)


// Test generated using Keploy
func TestGetURL_MissingHash(t *testing.T) {
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)

    c.Params = gin.Params{}
    getURL(c)

    if w.Code != http.StatusBadRequest {
        t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
    }
}

// Test generated using Keploy
func TestPutURL_InvalidRequestBody(t *testing.T) {
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)

    c.Request = httptest.NewRequest(http.MethodPost, "/", nil)
    putURL(c)

    if w.Code != http.StatusBadRequest {
        t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
    }
}


// Test generated using Keploy
func TestGenerateShortLink_ConsistentHash(t *testing.T) {
    input := "http://example.com"
    expectedLength := 8

    result := GenerateShortLink(input)

    if len(result) != expectedLength {
        t.Errorf("Expected hash length %d, got %d", expectedLength, len(result))
    }
}

