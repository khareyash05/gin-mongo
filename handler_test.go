package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"context"
	"time"

	"io/ioutil"
	"strings"

	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"go.uber.org/zap/zaptest"
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

// Test generated using Keploy
func TestGet_URLFound_123(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("should return URL when found in database", func(mt *mtest.T) {
		expectedURL := &URL{
			ID:      "test-id",
			Created: time.Now(),
			Updated: time.Now(),
			URL:     "http://example.com",
		}
		first := mtest.CreateCursorResponse(
			1,
			"testdb.testcoll",
			mtest.FirstBatch,
			bson.D{
				{Key: "_id", Value: expectedURL.ID},
				{Key: "created", Value: expectedURL.Created},
				{Key: "updated", Value: expectedURL.Updated},
				{Key: "url", Value: expectedURL.URL},
			},
		)
		second := mtest.CreateCursorResponse(0, "testdb.testcoll", mtest.NextBatch)
		mt.AddMockResponses(first, second)

		col = mt.Coll
		result, err := Get(context.TODO(), "test-id")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result.ID != expectedURL.ID || result.URL != expectedURL.URL {
			t.Errorf("expected %+v, got %+v", expectedURL, result)
		}
	})
}

// Test generated using Keploy
func TestUpsert_Success_456(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("should upsert URL successfully", func(mt *mtest.T) {
		col = mt.Coll
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		url := URL{
			ID:      "test-id",
			Created: time.Now(),
			Updated: time.Now(),
			URL:     "http://example.com",
		}

		err := Upsert(context.TODO(), url)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

// Test generated using Keploy
func TestPutURL_ValidRequestBody_321(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("should insert URL successfully", func(mt *mtest.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		col = mt.Coll
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		requestBody := `{"url": "http://example.com"}`
		c.Request = httptest.NewRequest(http.MethodPost, "/", nil)
		c.Request.Body = ioutil.NopCloser(strings.NewReader(requestBody))

		putURL(c)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}
	})
}

// Test generated using Keploy
func TestNew_Success_987(t *testing.T) {
	host := "localhost"
	db := "testdb"

	client, err := New(host, db)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if client == nil {
		t.Fatalf("expected a non-nil MongoDB client")
	}
}

// Test generated using Keploy
func TestUpsert_UpdateOneError_502(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("should return error when UpdateOne fails", func(mt *mtest.T) {
		col = mt.Coll
		// Mock UpdateOne to return an error
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   0,
			Code:    11000, // Example duplicate key error code
			Message: "Simulated UpdateOne error",
		}))

		urlToUpsert := URL{
			ID:      "test-id",
			Created: time.Now(),
			Updated: time.Now(),
			URL:     "http://example-error.com",
		}

		err := Upsert(context.TODO(), urlToUpsert)

		assert.Error(t, err)
	})
}

// Test generated using Keploy
func TestPutURL_UpsertError_505(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("should return 500 when Upsert fails", func(mt *mtest.T) {
		logger = zaptest.NewLogger(t) // Initialize logger
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		col = mt.Coll
		// Mock Upsert (UpdateOne) to return an error
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   0,
			Code:    12345, // Some error code
			Message: "Simulated Upsert error",
		}))

		requestBody := `{"url": "http://example-error.com"}`
		c.Request = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(requestBody))
		c.Request.Header.Set("Content-Type", "application/json")

		putURL(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		// Assert response body for error message structure
		var resp map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err, "Failed to unmarshal response body")
		assert.Contains(t, resp, "error", "Response should contain an error field")
		// Note: Asserting the exact error string might be brittle if it includes dynamic parts
	})
}
