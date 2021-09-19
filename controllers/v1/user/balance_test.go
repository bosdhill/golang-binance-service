package user

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/bosdhill/golang-binance-service/core/models"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var user models.User

func init() {
	gin.SetMode(gin.TestMode)

	err := godotenv.Load(".env.test")
	if err != nil {
		log.Fatal("Error loading .env.test file")
	}

	user.APIKey = os.Getenv("FUTURES_API_KEY")
	user.APISecret = os.Getenv("FUTURES_API_SECRET")
	futures.UseTestnet = true
}

func TestGetBalance(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	userJSON, err := json.Marshal(user)
	if err != nil {
		t.Fatal(err)
	}

	c.Request, err = http.NewRequest(
		"GET",
		"/foo",
		bytes.NewBuffer(userJSON),
	)

	if err != nil {
		t.Fatal(err)
	}

	GetBalance(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var got gin.H
	err = json.Unmarshal(w.Body.Bytes(), &got)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, got, "account empty")

	res, err := json.MarshalIndent(&got, "", " ")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf(string(res))
}
