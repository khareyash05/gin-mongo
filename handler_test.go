package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"go.uber.org/zap"
)

// Test generated using Keploy
func TestGet_ValidID_123(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("should retrieve URL document by ID", func(mt *mtest.T) {
		col = mt.Coll
		expectedURL := URL{
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
func TestUpsert_ValidURL_456(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("should upsert URL document", func(mt *mtest.T) {
		col = mt.Coll
		url := URL{
			ID:      "test-id",
			Created: time.Now(),
			Updated: time.Now(),
			URL:     "http://example.com",
		}

		mt.AddMockResponses(mtest.CreateSuccessResponse())

		err := Upsert(context.TODO(), url)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

// Test generated using Keploy
func TestGetURL_ValidHash_789(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("should redirect to URL when hash exists", func(mt *mtest.T) {
		col = mt.Coll
		logger = zap.NewNop()
		expectedURL := URL{
			ID:      "test-hash",
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

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "param", Value: "test-hash"}}
		c.Request = httptest.NewRequest(http.MethodGet, "/:param", nil)

		getURL(c)

		if w.Code != http.StatusSeeOther {
			t.Errorf("expected status %d, got %d", http.StatusSeeOther, w.Code)
		}
		if w.Header().Get("Location") != expectedURL.URL {
			t.Errorf("expected location %s, got %s", expectedURL.URL, w.Header().Get("Location"))
		}
	})
}

// Test generated using Keploy
func TestPutURL_ValidRequest_101(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("should generate short link and upsert URL", func(mt *mtest.T) {
		col = mt.Coll
		logger = zap.NewNop()
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		body := `{"url": "http://example.com"}`
		c.Request = httptest.NewRequest(http.MethodPut, "/", strings.NewReader(body))
		c.Request.Header.Set("Content-Type", "application/json")

		putURL(c)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}
	})
}

// Test generated using Keploy
func TestGet_NotFound_512(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("should return ErrNoDocuments when ID not found", func(mt *mtest.T) {
		col = mt.Coll
		mt.AddMockResponses(mtest.CreateCursorResponse(0, "testdb.testcoll", mtest.FirstBatch)) // No documents found

		_, err := Get(context.TODO(), "non-existent-id")
		require.Error(t, err)
		assert.Equal(t, mongo.ErrNoDocuments, err)
	})
}

// Test generated using Keploy
func TestGetURL_MissingParam_334(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	// No params set
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil) // Param would normally be in path

	getURL(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "please append url hash")
}

// Test generated using Keploy
func TestGetURL_GetError_NotFound_765(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("should return 404 when Get returns ErrNoDocuments", func(mt *mtest.T) {
		col = mt.Coll
		logger = zap.NewNop() // Use Noop logger for testing

		mt.AddMockResponses(mtest.CreateCursorResponse(0, "testdb.testcoll", mtest.FirstBatch)) // Simulate not found

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{gin.Param{Key: "param", Value: "not-found-hash"}}
		c.Request, _ = http.NewRequestWithContext(context.TODO(), http.MethodGet, "/not-found-hash", nil)

		getURL(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "url not found")
	})
}

// Test generated using Keploy
func TestPutURL_BindError_281(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	// Malformed JSON
	body := `{"url": "http://example.com"`
	c.Request = httptest.NewRequest(http.MethodPut, "/", strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	putURL(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "failed to decode req")
}

// Test generated using Keploy
func TestPutURL_MissingURLField_804(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	// JSON missing the 'url' field
	body := `{"other_field": "some_value"}`
	c.Request = httptest.NewRequest(http.MethodPut, "/", strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	putURL(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "missing url param")
}

// Test generated using Keploy
func TestNew_ValidHostAndDB_904(t *testing.T) {
	host := "localhost"
	db := "testdb"

	client, err := New(host, db)

	assert.NoError(t, err)
	assert.NotNil(t, client)
}
