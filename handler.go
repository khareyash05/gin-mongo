package main

import (
	"context"
	"crypto/sha256"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/itchyny/base58-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type URL struct {
	ID      string    `json:"id" bson:"_id"`
	Created time.Time `json:"created" bson:"created"`
	Updated time.Time `json:"updated" bson:"updated"`
	URL     string    `json:"URL" bson:"url"`
}

func Get(ctx context.Context, id string) (*URL, error) {
	filter := bson.M{"_id": id}
	var u URL
	err := col.FindOne(ctx, filter).Decode(&u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func putURL(c *gin.Context) {
	var m map[string]string

	err := c.ShouldBindJSON(&m)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to decode req"})
		return
	}
	u := m["url"]

	if u == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing url param"})
		return
	}

	id := GenerateShortLink(u)
	c.JSON(http.StatusOK, gin.H{
		"ts":  time.Now().UnixNano(),
		"url": "http://localhost:8080/" + id,
	})
}

func New(host, db string) (*mongo.Client, error) {
	clientOptions := options.Client()

	clientOptions.ApplyURI("mongodb://" + host + "/" + db + "?retryWrites=true&w=majority")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return mongo.Connect(ctx, clientOptions)
}

func GenerateShortLink(initialLink string) string {
	urlHashBytes := sha256Of(initialLink)
	generatedNumber := new(big.Int).SetBytes(urlHashBytes).Uint64()
	finalString := base58Encoded([]byte(fmt.Sprintf("%d", generatedNumber)))
	return finalString[:8]
}

func sha256Of(input string) []byte {
	algorithm := sha256.New()
	algorithm.Write([]byte(input))
	return algorithm.Sum(nil)
}

func base58Encoded(bytes []byte) string {
	encoding := base58.BitcoinEncoding
	encoded, _ := encoding.Encode(bytes)
	return string(encoded)
}
