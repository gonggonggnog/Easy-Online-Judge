package middleware

import (
	"blog/helper"
	"github.com/gin-gonic/gin"
)

func AuthAdminCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization") // 获取请求头中的Authorization
		userClaim, err := helper.ParseToken(auth)
		if err != nil || userClaim == nil {
			c.JSON(200, gin.H{
				"code": -1,
				"msg":  "unauthorized Authorization",
			})
			c.Abort()
			return
		}
		if userClaim.IsAdmin != 1 {
			c.JSON(200, gin.H{
				"code": -1,
				"msg":  "unauthorized Admin",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
