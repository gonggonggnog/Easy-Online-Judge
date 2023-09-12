package service

import (
	"blog/dao"
	"blog/helper"
	"blog/models"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"math/rand"
	"strconv"
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

// SendCode
// @Tags 公共方法
// @Summary 发送验证码
// @Param email formData string true "email"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /send-code [post]
func SendCode(c *gin.Context) {
	email := c.PostForm("email")
	if email == "" {
		c.JSON(200, gin.H{
			"code": -1,
			"data": "参数不正确",
		})
	}
	code := strconv.Itoa(rand.Int()%1000000 + 100000)
	dao.RedisSet(email, code)
	err := helper.SendEmail("405351435@qq.com", code)
	if err != nil {
		c.JSON(200, gin.H{
			"code": -1,
			"data": "send email failed:" + err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"data": "send email success",
	})
}

// Register
// @Tags 公共方法
// @Summary 注册用户
// @Param email formData string true "email"
// @Param code formData string true "code"
// @Param name formData string true "name"
// @Param password formData string true "password"
// @Param phone formData string false "phone"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /register [post]

func Register(c *gin.Context) {
	email := c.PostForm("email")
	userCode := c.PostForm("code")
	name := c.PostForm("name")
	password := c.PostForm("password")
	phone := c.PostForm("phone")
	if email == "" || userCode == "" || name == "" || password == "" {
		c.JSON(200, gin.H{
			"code": -1,
			"data": "必填项目不能为空",
		})
		return
	}
	sysCode, err := dao.RedisGet(email)
	if err != nil {
		c.JSON(200, gin.H{
			"code": -1,
			"data": "redis get code error:" + err.Error(),
		})
		return
	}
	if sysCode != userCode {
		c.JSON(200, gin.H{
			"code": -1,
			"data": "验证码不正确，请重新获取验证码",
		})
		return
	}
	UserIdentity := helper.GenerateUUid()
	data := &models.UserBasic{
		Identity: UserIdentity,
		Name:     name,
		Password: helper.GetMd5(password),
		Email:    email,
		Phone:    phone,
	}
	err = dao.DB.Create(data).Error
	if err != nil {
		c.JSON(200, gin.H{
			"code": -1,
			"data": "user create error" + err.Error(),
		})
		return
	}
	var token string
	token, err = helper.GenerateToken(UserIdentity, name, 0)
	if err != nil {
		c.JSON(200, gin.H{
			"code": -1,
			"data": "token generate error:" + err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"data": gin.H{
			"token": token,
		},
	})
}
