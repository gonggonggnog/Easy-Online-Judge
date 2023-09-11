package service

import (
	"blog/dao"
	"blog/helper"
	"blog/models"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetUserDetail
// @Tags 公共方法
// @Summary 用户详情
// @Param identity query string false "user identity"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /user-detail [get]
func GetUserDetail(c *gin.Context) {
	identity := c.Query("identity")
	if identity == "" {
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  "用户标识不能为空",
		})
		return
	}
	data := new(models.UserBasic)
	err := dao.DB.Where("identity=?", identity).Omit("password").First(&data).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(200, gin.H{
				"code": -1,
				"msg":  "用户标识不存在",
			})
			return
		}
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  "Get ProblemDetail Error:" + err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  data,
	})
}

// Login
// @Tags 公共方法
// @Summary 用户登录
// @Param username formData string false "username"
// @Param password formData string false "password"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /login [post]
func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	if username == "" || password == "" {
		c.JSON(200, gin.H{
			"code": -1,
			"data": "必填信息为空",
		})
		return
	}
	data, err := dao.GetPasswd(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(200, gin.H{
				"code": -1,
				"data": "账号不存在",
			})
			return
		}
		c.JSON(200, gin.H{
			"code": -1,
			"data": "err occur while get passwd:" + err.Error(),
		})
		return
	}
	//fmt.Println(helper.GetMd5("123123"))
	if data.Password != helper.GetMd5(password) {
		c.JSON(200, gin.H{
			"code": -1,
			"data": "账号密码不匹配",
		})
		return
	}
	tokenString, err := helper.GenerateToken(data.Identity, data.Name, data.IsAdmin)
	if err != nil {
		c.JSON(200, gin.H{
			"code": -1,
			"data": "generate toke  failed:" + err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"data": gin.H{
			"msg":   "登录成功",
			"token": tokenString,
		},
	})
}
