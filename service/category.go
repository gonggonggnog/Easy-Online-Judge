package service

import (
	"blog/dao"
	"blog/define"
	"blog/helper"
	"blog/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

// GetCategoryList
// @Tags 管理员私有方法
// @Summary 分类列表
// @Param Authorization header string true "token"
// @Param page query int false "page"
// @Param size query int false "size"
// @Param keyword query string false "keyword"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /admin/category-list [get]
func GetCategoryList(c *gin.Context) {
	size, _ := strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))
	page, _ := strconv.Atoi(c.DefaultQuery("page", define.DefaultPage))
	page = (page - 1) * size
	keyword := c.Query("keyword")
	tx := dao.GetCategoryList(keyword)
	var count int64
	list := make([]*models.CategoryBasic, 0)
	err := tx.Count(&count).Offset(page).Limit(size).Find(&list).Error // 分页查询
	if err != nil {
		c.JSON(200, gin.H{
			"code": -1,
			"data": "select categories error:" + err.Error(),
		})
	}

	c.JSON(200, gin.H{
		"code": 200,
		"data": gin.H{
			"count": count,
			"list":  list,
		},
	})
}

// AddCategory
// @Tags 管理员私有方法
// @Summary 添加分类
// @Param Authorization header string true "token"
// @Param name formData string true "name"
// @Param parent_id formData string false "parent_id"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /admin/category-create [post]
func AddCategory(c *gin.Context) {
	name := c.PostForm("name")
	parentId, _ := strconv.Atoi(c.PostForm("parent_id"))
	if name == "" {
		c.JSON(200, gin.H{
			"code": -1,
			"data": "必填字段为空",
		})
		return
	}
	var count int64
	dao.DB.Model(&models.CategoryBasic{}).Where("name=?", name).Count(&count) // 判断分类是否存在
	if count > 0 {
		c.JSON(200, gin.H{
			"code": -1,
			"data": "分类已存在",
		})
		return
	}
	err := dao.DB.Create(&models.CategoryBasic{
		Identity: helper.GenerateUUid(),
		Name:     name,
		ParentId: uint(parentId),
	}).Error
	if err != nil {
		c.JSON(200, gin.H{
			"code": -1,
			"data": "create category error:" + err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"data": "create category success",
	})
}

// UpdateCategory
// @Tags 管理员私有方法
// @Summary 更新分类
// @Param Authorization header string true "token"
// @Param identity formData string true "identity"
// @Param name formData string false "name"
// @Param parent_id formData string false "parent_id"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /admin/category-update [put]
func UpdateCategory(c *gin.Context) {
	identity := c.PostForm("identity")
	name := c.PostForm("name")
	parentId, _ := strconv.Atoi(c.PostForm("parent_id"))
	if name == "" {
		c.JSON(200, gin.H{
			"code": -1,
			"data": "必填字段为空",
		})
		return
	}
	err := dao.DB.Model(&models.CategoryBasic{}).Where("identity=?", identity).Updates(&models.CategoryBasic{ // 判断分类是否存在
		Name:     name,
		ParentId: uint(parentId),
	}).Error
	if err != nil {
		c.JSON(200, gin.H{
			"code": -1,
			"data": "update category error:" + err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"data": "update category success",
	})
}

// DeleteCategory
// @Tags 管理员私有方法
// @Summary 删除分类
// @Param Authorization header string true "token"
// @Param identity query string true "identity"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /admin/category-delete [delete]
func DeleteCategory(c *gin.Context) {
	identity := c.Query("identity")
	fmt.Println(identity)
	if identity == "" {
		c.JSON(200, gin.H{
			"code": -1,
			"data": "必填字段为空",
		})
		return
	}
	var count int64
	dao.DB.Model(new(models.ProblemsCategory)).
		Where("category_id=(select id from category_basic where category_basic.identity=? LIMIT 1)", identity).
		Count(&count) // 判断分类下是否存在问题
	if count > 0 {
		c.JSON(200, gin.H{
			"code": -1,
			"data": "该分类下存在问题，无法删除",
		})
		return
	}
	err := dao.DB.Where("identity=?", identity).Delete(&models.CategoryBasic{}).Error // 删除分类
	if err != nil {
		c.JSON(200, gin.H{
			"code": -1,
			"data": "delete category error:" + err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"data": "delete category success",
	})
}
