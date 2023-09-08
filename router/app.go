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
	r.GET("/problem-list", service.GetProblemList)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.GET("/problem-detail", service.GetProblemDetail)
	r.GET("/user-detail", service.GetUserDetail)
	r.GET("/submit-list", service.GetSubmitList)
	return r
}
