package middleware

import (
	"github.com/gin-gonic/gin"
	"go-google-cloud-storage/app/constant"
	"go-google-cloud-storage/app/helper"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiID := helper.GenerateApiCallID()
		c.Set(constant.RequestIDKey, apiID)

		c.Next()
	}
}
