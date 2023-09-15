package service

import (
	"blog/dao"
	"blog/define"
	"blog/helper"
	"blog/models"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
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
	err := dao.DB.Where("identity=?", identity).Omit("password").First(&data).Error // Omit("password")忽略password字段
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { // 判断是否为记录未找到的错误
			c.JSON(200, gin.H{
				"code": -1,
				"msg":  "用户标识不存在",
			})
			return
		}
		c.JSON(200, gin.H{ // 其他错误
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
	if username == "" || password == "" { // 判断是否为空
		c.JSON(200, gin.H{
			"code": -1,
			"data": "必填信息为空",
		})
		return
	}
	data, err := dao.GetPasswd(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { // 判断是否为记录未找到的错误
			c.JSON(200, gin.H{
				"code": -1,
				"data": "账号不存在",
			})
			return
		}
		c.JSON(200, gin.H{ // 其他错误
			"code": -1,
			"data": "err occur while get passwd:" + err.Error(),
		})
		return
	}
	if data.Password != helper.GetMd5(password) { // 判断密码是否正确
		c.JSON(200, gin.H{
			"code": -1,
			"data": "账号密码不匹配",
		})
		return
	}
	tokenString, err := helper.GenerateToken(data.Identity, data.Name, data.IsAdmin) // 生成token
	if err != nil {
		c.JSON(200, gin.H{
			"code": -1,
			"data": "generate toke failed:" + err.Error(),
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
	code := strconv.Itoa(rand.Int()%899999 + 100000)
	dao.RedisSet(email, code)                         // 将验证码存入redis
	err := helper.SendEmail("405351435@qq.com", code) // 发送邮件
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
// @Summary 用户注册
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
	var isExists int64
	dao.DB.Model(new(models.UserBasic)).Where("email = ?", email).Count(&isExists) // 判断邮箱是否已被注册
	if isExists > 0 {
		c.JSON(200, gin.H{
			"code": -1,
			"data": "该邮箱已被注册",
		})
		return
	}
	sysCode, err := dao.RedisGet(email)
	if sysCode != userCode || err != nil {
		c.JSON(200, gin.H{
			"code": -1,
			"data": "验证码不正确，请重新获取验证码",
		})
		log.Println(err.Error())
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

// GetRankList
// @Tags 公共方法
// @Summary 用户排行榜
// @Param page query int false "page"
// @Param size query int false "size"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /rank-list [get]
func GetRankList(c *gin.Context) {
	var count int64
	list := make([]*models.UserBasic, 0)
	page, err := strconv.Atoi(c.DefaultQuery("page", define.DefaultPage))
	if err != nil {
		log.Fatal("get problem page error:", err)
		return
	}
	size, terr := strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))
	if terr != nil {
		log.Fatal("get problem size error:", err)
		return
	}
	page = (page - 1) * size
	err = dao.DB.Model(new(models.UserBasic)).Count(&count).Select("id,name,pass_num,submit_num").
		Order("pass_num DESC,submit_num ASC").
		Offset(page).Limit(size).Find(&list).Error // 分页查询
	if err != nil {
		c.JSON(200, gin.H{
			"code": -1,
			"data": "get rank list error:" + err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"data": gin.H{
			"list":  list,
			"count": count,
		},
	})
}
