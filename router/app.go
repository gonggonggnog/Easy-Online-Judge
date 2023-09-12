package router

import (
	_ "blog/docs"
	"blog/service"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(context *gin.Context) {
		context.String(200, "pong")
	})
	r.GET("/problem-list", service.GetProblemList) //获取问题列表
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.GET("/problem-detail", service.GetProblemDetail) //获取问题详情
	r.GET("/user-detail", service.GetUserDetail)       //获取用户详情
	r.GET("/submit-list", service.GetSubmitList)       //获取提交列表
	r.POST("/login", service.Login)
	r.POST("/send-code", service.SendCode)
	r.POST("/register", service.Register)
	return r
}
