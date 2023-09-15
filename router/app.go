package router

import (
	_ "blog/docs"
	"blog/middleware"
	"blog/service"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() *gin.Engine {
	r := gin.Default()
	//swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler)) // default ginSwagger url
	//公共方法
	//问题接口
	r.GET("/problem-list", service.GetProblemList)     //获取问题列表
	r.GET("/rank-list", service.GetRankList)           //用户排行榜
	r.GET("/problem-detail", service.GetProblemDetail) //获取问题详情
	//用户接口
	r.GET("/user-detail", service.GetUserDetail) //获取用户详情
	r.POST("/login", service.Login)              //用户登录
	r.POST("/register", service.Register)        //用户注册
	//其他接口
	r.GET("/submit-list", service.GetSubmitList) //获取提交列表
	r.POST("/send-code", service.SendCode)       //发送验证码

	//管理员方法
	Admin := r.Group("/admin", middleware.AuthAdminCheck())
	//问题接口
	Admin.POST("/problem-create", service.AddProblem)    //添加问题
	Admin.PUT("/problem-update", service.UpdateProblem)  //更新问题
	Admin.POST("/problem-delete", service.DeleteProblem) //删除问题
	//分类接口
	Admin.GET("/category-list", service.GetCategoryList)     //获取分类列表
	Admin.POST("/category-create", service.AddCategory)      //添加分类
	Admin.PUT("/category-update", service.UpdateCategory)    //更新分类
	Admin.DELETE("/category-delete", service.DeleteCategory) //删除分类

	//用户接口
	User := r.Group("/user", middleware.AuthUserCheck()) //用户登录后才能访问
	//提交接口
	User.POST("/submit", service.Submit) //添加提交
	return r
}
