package middleware

import (
	"net/http"

	"github.com/bosdhill/golang-binance-service/core/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func validatorFunc(c *gin.Context, obj interface{}) {
	if err := c.ShouldBindJSON(&obj); err == nil {
		validate := validator.New()
		if err := validate.Struct(&obj); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			c.Abort()
			return
		}
	}
	c.Next()
}

func CreateOrderValidator(c *gin.Context) {
	var order models.TakeProfitOrder
	validatorFunc(c, order)
}

func GetBalanceValidator(c *gin.Context) {
	var user models.User
	validatorFunc(c, user)
}
