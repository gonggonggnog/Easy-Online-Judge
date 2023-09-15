package middleware

import (
	"blog/helper"
	"github.com/gin-gonic/gin"
)

func AuthUserCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		userClaim, err := helper.ParseToken(auth)
		if err != nil || userClaim == nil {
			c.JSON(200, gin.H{
				"code": -1,
				"msg":  "unauthorized Authorization",
			})
			c.Abort()
			return
		}
		c.Set("user", userClaim)
		c.Next()

	}
}
