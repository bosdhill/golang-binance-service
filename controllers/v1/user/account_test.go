package user

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/bosdhill/golang-binance-service/core/models"
	"github.com/bosdhill/golang-binance-service/libs/test"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	test.IntializeControllerTests()
}

func TestGetAccount(t *testing.T) {
	var user models.User
	user.APIKey = os.Getenv("FUTURES_API_KEY")
	user.APISecret = os.Getenv("FUTURES_API_SECRET")

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

	GetAccount(c)

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
