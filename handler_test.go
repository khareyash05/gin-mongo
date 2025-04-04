package main

import (
	"testing"

	"context"
	"time"

	"net/http"
	"net/http/httptest"

	"strings"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"go.uber.org/zap"
)

// Test generated using Keploy
func TestPutURL_ValidRequest_321(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("should create short link and store in database", func(mt *mtest.T) {
		col = mt.Coll
		logger = zap.NewNop()
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		body := `{"url": "http://example.com"}`
		c.Request = httptest.NewRequest(http.MethodPut, "/put-url", strings.NewReader(body))
		c.Request.Header.Set("Content-Type", "application/json")

		putURL(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "http://localhost:8080/")
	})
}

// Test generated using Keploy

// Test generated using Keploy


func TestGetURL_ValidHash_789(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("should redirect to URL for valid hash", func(mt *mtest.T) {
		col = mt.Coll
		logger = zap.NewNop()
		expectedURL := URL{
			ID:      "test-hash",
			Created: time.Now(),
			Updated: time.Now(),
			URL:     "http://example.com",
		}

		mt.AddMockResponses(mtest.CreateCursorResponse(
			1,
			"testdb.testcoll",
			mtest.FirstBatch,
			bson.D{
				{Key: "_id", Value: expectedURL.ID},
				{Key: "created", Value: expectedURL.Created},
				{Key: "updated", Value: expectedURL.Updated},
				{Key: "url", Value: expectedURL.URL},
			},
		), mtest.CreateCursorResponse(0, "testdb.testcoll", mtest.NextBatch))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "param", Value: "test-hash"}}
		c.Request = httptest.NewRequest(http.MethodGet, "/test-hash", nil)

		getURL(c)

		assert.Equal(t, http.StatusSeeOther, w.Code)
		assert.Equal(t, "http://example.com", w.Header().Get("Location"))
	})
}

// Test generated using Keploy

// Test generated using Keploy


func TestPutURL_BadJSON_578(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	body := `{"url": "http://example.com",}` // Invalid JSON (trailing comma)
	c.Request = httptest.NewRequest(http.MethodPut, "/put-url", strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	putURL(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `{"error": "failed to decode req"}`, w.Body.String())
}

// Test generated using Keploy


func TestPutURL_MissingURLParam_601(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	body := `{"noturl": "http://example.com"}` // Missing 'url' field
	c.Request = httptest.NewRequest(http.MethodPut, "/put-url", strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	putURL(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `{"error": "missing url param"}`, w.Body.String())
}


func TestGetURL_MissingHash_115(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	// No params set, simulating missing ":param"

	getURL(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `{"error": "please append url hash"}`, w.Body.String())
}

// Test generated using Keploy


func TestUpsert_DBError_991(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("should return error on db failure", func(mt *mtest.T) {
		col = mt.Coll // Set the global collection
		url := URL{
			ID:      "test-error-id",
			Created: time.Now(),
			Updated: time.Now(),
			URL:     "http://error.example.com",
		}

		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   0,
			Code:    11000, // Example error code (duplicate key)
			Message: "db write error",
		}))

		err := Upsert(context.TODO(), url)
		require.Error(t, err)
		// Check if the error is a mongo write exception if necessary
		// assert.True(t, mongo.IsDuplicateKeyError(err))
	})
}

// Test generated using Keploy

// Test generated using Keploy

func TestNew_ValidInitialization_654(t *testing.T) {
	host := "localhost"
	db := "testdb"

	client, err := New(host, db)

	assert.NoError(t, err, "New should not return an error")
	assert.NotNil(t, client, "MongoDB client should not be nil")
}

