package service

import (
	"blog/dao"
	"blog/define"
	"blog/models"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"strconv"
)

// GetProblemList
// @Tags 公共方法
// @Summary 问题列表
// @Param page query int false "page"
// @Param size query int false "size"
// @Param keyword query string false "keyword"
// @Param category_identity query string false "category_identity"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /problem-list [get]
func GetProblemList(c *gin.Context) {
	var count int64
	list := make([]*models.ProblemsBasic, 0)
	page, err := strconv.Atoi(c.DefaultQuery("page", define.DefaultPage))
	if err != nil {
		log.Fatal("get problem page error:", err)
		return
	}
	categoryIdentity := c.Query("category_identity")
	size, terr := strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))
	if terr != nil {
		log.Fatal("get problem size error:", err)
		return
	}
	keyword := c.Query("keyword")
	page = (page - 1) * size
	tx := dao.GetProblemList(keyword, categoryIdentity)
	err = tx.Count(&count).Omit("content").Offset(page).Limit(size).Find(&list).Error
	if err != nil {
		log.Fatal("get problem error:", err)
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

// GetProblemDetail
// @Tags 公共方法
// @Summary 问题详情
// @Param identity query string false "problem identity"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /problem-detail [get]
func GetProblemDetail(c *gin.Context) {
	identity := c.Query("identity")
	if identity == "" {
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  "问题标识不能为空",
		})
		return
	}
	data := new(models.ProblemsBasic)
	err := dao.DB.Preload("ProblemCategories").Preload("ProblemCategories.CategoryBasic").Where("identity=?", identity).First(&data).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(200, gin.H{
				"code": -1,
				"msg":  "问题不存在",
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
