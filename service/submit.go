package service

import (
	"blog/dao"
	"blog/define"
	"blog/models"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
)

// GetSubmitList
// @Tags 公共方法
// @Summary 用户详情
// @Param page query int false "page"
// @Param size query int false "size"
// @Param problem_identity query string false "problem identity"
// @Param user_identity query string false "user identity"
// @Param status query int false "status"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /submit-list [get]
func GetSubmitList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", define.DefaultPage))
	size, _ := strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))
	//status, _ := strconv.Atoi(c.DefaultQuery("status", define.DefaultStatus))
	problemIdentity := c.Query("problem_identity")
	userIdentity := c.Query("user_identity")
	data := new(models.SubmitBasic)
	page = (page - 1) * size
	err := dao.DB.Where("user_identity=? and problem_identity=?", userIdentity, problemIdentity).Offset(page).Limit(size).Find(&data).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(200, gin.H{
				"code": -1,
				"msg":  "提交数据不存在",
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
