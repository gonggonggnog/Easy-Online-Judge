package service

import (
	"blog/dao"
	"blog/define"
	"blog/helper"
	"blog/models"
	"encoding/json"
	"errors"
	"fmt"
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
	err := dao.DB.Preload("ProblemCategories").Preload("ProblemCategories.CategoryBasic").Preload("TestCase").Where("identity=?", identity).First(&data).Error
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

// AddProblem
// @Tags 管理员私有方法
// @Summary 问题创建
// @param Authorization header string true "Authorization"
// @Param title formData string true "title"
// @Param content formData string true "content"
// @Param max_runtime formData int false "max_runtime"
// @Param max_mem formData int false "max_mem"
// @Param category_ids formData []string false "category_ids" collectionFormat(multi)
// @Param test_case formData []string true "test_case" collectionFormat(multi)
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /admin/problem-create [post]
func AddProblem(c *gin.Context) {
	title := c.PostForm("title")
	content := c.PostForm("content")
	maxRuntime, _ := strconv.Atoi(c.PostForm("max_runtime"))
	maxMem, _ := strconv.Atoi(c.PostForm("max_mem"))
	categoryIds := c.PostFormArray("category_ids")
	testCase := c.PostFormArray("test_case")
	if title == "" || content == "" || len(testCase) == 0 {
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  "必填信息不能为空",
		})
		return
	}
	// 1. 创建问题基本信息
	data := &models.ProblemsBasic{
		Identity:   helper.GenerateUUid(),
		Title:      title,
		Content:    content,
		MaxRuntime: maxRuntime,
		MaxMem:     maxMem,
	}
	// 2. 创建问题分类信息
	problemsCategories := make([]*models.ProblemsCategory, 0)

	for _, v := range categoryIds {
		categoryId, _ := strconv.Atoi(v)
		problemsCategories = append(problemsCategories, &models.ProblemsCategory{
			ProblemId:  data.ID,
			CategoryId: uint(categoryId),
		})
	}
	data.ProblemCategories = problemsCategories
	// 3. 创建问题测试用例
	testCases := make([]*models.TestCase, 0)
	for _, v := range testCase {
		testCaseMap := make(map[string]string)
		err := json.Unmarshal([]byte(v), &testCaseMap)
		if err != nil {
			c.JSON(200, gin.H{
				"code": -1,
				"msg":  "测试样例格式错误",
			})
			return
		}
		// {"input":"1 2","output":"3"}
		if _, ok := testCaseMap["input"]; !ok {
			c.JSON(200, gin.H{
				"code": -1,
				"msg":  "测试样例格式错误",
			})
			return
		}
		if _, ok := testCaseMap["output"]; !ok {
			c.JSON(200, gin.H{
				"code": -1,
				"msg":  "测试样例格式错误",
			})
			return
		}
		testCases = append(testCases, &models.TestCase{
			Identity:        helper.GenerateUUid(),
			ProblemIdentity: data.Identity,
			Input:           testCaseMap["input"],
			Output:          testCaseMap["output"],
		})
	}
	data.TestCase = testCases
	err := dao.DB.Create(&data).Error // 创建问题基本信息
	if err != nil {
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  "Create ProblemBasic Error:" + err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "创建成功",
	})

}

// UpdateProblem
// @Tags 管理员私有方法
// @Summary 问题更新
// @param Authorization header string true "Authorization"
// @Param identity formData string true "identity"
// @Param title formData string false "title"
// @Param content formData string false "content"
// @Param max_runtime formData int false "max_runtime"
// @Param max_mem formData int false "max_mem"
// @Param category_ids formData []string false "category_ids" collectionFormat(multi)
// @Param test_case formData []string false "test_case" collectionFormat(multi)
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /admin/problem-update [put]
func UpdateProblem(c *gin.Context) {
	title := c.PostForm("title")
	content := c.PostForm("content")
	maxRuntime, _ := strconv.Atoi(c.PostForm("max_runtime"))
	maxMem, _ := strconv.Atoi(c.PostForm("max_mem"))
	categoryIds := c.PostFormArray("category_ids")
	testCase := c.PostFormArray("test_case")
	identity := c.PostForm("identity")
	if identity == "" {
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  "问题标识不能为空",
		})
		return
	}
	if err := dao.DB.Transaction(func(tx *gorm.DB) error {
		problemBasic := &models.ProblemsBasic{
			Identity:   identity,
			Title:      title,
			Content:    content,
			MaxRuntime: maxRuntime,
			MaxMem:     maxMem,
		}
		err := tx.Model(&models.ProblemsBasic{}).Where("identity=?", identity).Updates(problemBasic).Error // 更新问题基本信息
		if err != nil {
			return err
		}
		err = tx.Where("identity=?", identity).Find(problemBasic).Error // 查询问题基本信息
		if err != nil {
			return err
		}
		// 1. 更新问题分类信息
		// 1.1 删除原有的分类信息
		err = tx.Where("problem_id=?", problemBasic.ID).Delete(&models.ProblemsCategory{}).Error // 删除问题分类信息
		if err != nil {
			return err
		}
		// 1.2 创建新的分类信息
		problemsCategories := make([]*models.ProblemsCategory, 0)

		for _, v := range categoryIds {
			categoryId, _ := strconv.Atoi(v)
			problemsCategories = append(problemsCategories, &models.ProblemsCategory{
				ProblemId:  problemBasic.ID,
				CategoryId: uint(categoryId),
			})
		}
		for _, v := range problemsCategories {
			fmt.Println(v)
		}
		err = tx.Create(&problemsCategories).Error
		if err != nil {
			return err
		}
		// 2. 更新问题测试用例
		// 2.1 删除原有的测试用例
		err = tx.Where("problem_identity=?", identity).Delete(&models.TestCase{}).Error
		if err != nil {
			return err
		}
		// 2.2 创建新的测试用例
		testCases := make([]*models.TestCase, 0)
		for _, v := range testCase {
			testCaseMap := make(map[string]string)
			err := json.Unmarshal([]byte(v), &testCaseMap)
			if err != nil {
				return err
			}
			// {"input":"1 2\n","output":"3\n"}
			if _, ok := testCaseMap["input"]; !ok {
				return errors.New("测试样例格式错误")
			}
			if _, ok := testCaseMap["output"]; !ok {
				return errors.New("测试样例格式错误")
			}
			testCases = append(testCases, &models.TestCase{
				Identity:        helper.GenerateUUid(),
				ProblemIdentity: problemBasic.Identity,
				Input:           testCaseMap["input"],
				Output:          testCaseMap["output"],
			})
		}
		err = tx.Create(&testCases).Error
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  "Update Problem Error:" + err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "更新成功",
	})
}

// DeleteProblem
// @Tags 管理员私有方法
// @Summary 问题删除
// @param Authorization header string true "Authorization"
// @Param identity formData string true "identity"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /admin/problem-delete [delete]
func DeleteProblem(c *gin.Context) {
	identity := c.PostForm("identity")
	if identity == "" {
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  "问题标识不能为空",
		})
		return
	}
	err := dao.DB.Transaction(func(tx *gorm.DB) error {
		// 1. 删除问题基本信息
		err := tx.Where("identity=?", identity).Delete(&models.ProblemsBasic{}).Error // 删除问题基本信息
		if err != nil {
			return err
		}
		// 2. 删除问题分类信息
		err = tx.Where("problem_identity=?", identity).Delete(&models.ProblemsCategory{}).Error // 删除问题分类信息
		if err != nil {
			return err
		}
		// 3. 删除问题测试用例
		err = tx.Where("problem_identity=?", identity).Delete(&models.TestCase{}).Error // 删除问题测试用例
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  "Delete Problem Error:" + err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "删除成功",
	})

}
