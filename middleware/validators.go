package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func Validator(c *gin.Context) {
	obj := c.Request.Body
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
