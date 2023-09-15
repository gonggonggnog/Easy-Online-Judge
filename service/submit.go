package service

import (
	"blog/dao"
	"blog/define"
	"blog/helper"
	"blog/models"
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io"
	"log"
	"os/exec"
	"runtime"
	"strconv"
	"sync"
	"time"
)

// GetSubmitList
// @Tags 公共方法
// @Summary 提交列表
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
	status, _ := strconv.Atoi(c.DefaultQuery("status", define.DefaultStatus))
	problemIdentity := c.Query("problem_identity")
	userIdentity := c.Query("user_identity")
	list := make([]*models.SubmitBasic, 0)
	page = (page - 1) * size
	var count int64
	tx := dao.GetSubmitList(userIdentity, problemIdentity, status)     // 获取提交列表
	err := tx.Count(&count).Offset(page).Limit(size).Find(&list).Error // 分页查询
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
			"msg":  "Get SubmitDetail Error:" + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"data": gin.H{
			"count": count,
			"list":  list,
		},
	})
}

// Submit
// @Tags 用户私有方法
// @Summary 代码提交
// @Param Authorization header string true "token"
// @Param problem_identity query string true "problem identity"
// @Param code body string true "code"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /user/submit [post]
func Submit(c *gin.Context) {
	problemIdentity := c.Query("problem_identity")
	code, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  "read code error:" + err.Error(),
		})
		return
	}
	if problemIdentity == "" {
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  "problem identity is empty",
		})
		return
	}
	file, err := helper.SaveCodeToFile(code)
	if err != nil {
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  "save code to file error:" + err.Error(),
		})
		return
	}

	userClaim, exist := c.Get("user")
	if !exist {
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  "user not exist",
		})
		return
	}
	userIdentity := userClaim.(*helper.UserClaims).Identity
	submit := &models.SubmitBasic{
		Identity:        helper.GenerateUUid(),
		UserIdentity:    userIdentity,
		ProblemIdentity: problemIdentity,
		Path:            file,
	}
	problemBasic := new(models.ProblemsBasic)
	err = dao.DB.Where("identity=?", problemIdentity).Preload("TestCase").First(problemBasic).Error // 获取问题基本信息
	if err != nil {
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  "problem not exist",
		})
		return
	}
	// 代码判断
	var msg string
	WA := make(chan int)
	OOM := make(chan int)
	CE := make(chan int)
	AC := make(chan int)
	passCount := 0
	var lock sync.Mutex
	for _, testCase := range problemBasic.TestCase {
		testCase := testCase
		go func() {
			command := exec.Command("go", "run", file)
			var out, stderr bytes.Buffer
			command.Stdout = &out
			command.Stderr = &stderr
			stinPipe, err := command.StdinPipe() // 获取标准输入
			if err != nil {
				log.Fatalln(err)
			}
			io.WriteString(stinPipe, testCase.Input)
			var bm runtime.MemStats
			runtime.ReadMemStats(&bm) // 获取内存
			err = command.Run()
			var em runtime.MemStats
			runtime.ReadMemStats(&em)
			if err != nil {
				msg = stderr.String()
				CE <- 1
				return
			}
			//答案错误
			if testCase.Output != out.String() {
				WA <- 1
				return
			}
			//超内存
			if em.Alloc-bm.Alloc > uint64(problemBasic.MaxMem)*1024 { // 超内存判断
				OOM <- 1
				return
			}
			lock.Lock()
			passCount++
			lock.Unlock()
			if passCount == len(problemBasic.TestCase) {
				AC <- 1
			}
		}()
	}
	dao.DB.Model(new(models.ProblemsBasic)).Where("identity=?", problemIdentity).
		Update("submit_num", gorm.Expr("submit_num+?", 1)) // 更新提交次数
	dao.DB.Model(new(models.UserBasic)).Where("identity=?", userIdentity).
		Update("submit_num", gorm.Expr("submit_num+?", 1)) // 更新提交次数
	select {
	case <-WA:
		msg = "答案错误"
		submit.Status = 2
	case <-OOM:
		msg = "运行超内存"
		submit.Status = 4
	case <-CE:
		submit.Status = 5
	case <-AC:
		msg = "答案正确"
		dao.DB.Model(new(models.ProblemsBasic)).Where("identity=?", problemIdentity).
			Update("pass_num", gorm.Expr("pass_num+?", 1)) // 更新通过次数
		dao.DB.Model(new(models.UserBasic)).Where("identity=?", userIdentity).
			Update("pass_num", gorm.Expr("pass_num+?", 1)) // 更新通过次数
		submit.Status = 1
	case <-time.After(time.Millisecond * time.Duration(problemBasic.MaxRuntime)): // 超时判断
		if passCount == len(problemBasic.TestCase) {
			msg = "答案正确"
			submit.Status = define.Accepted
		} else {
			msg = "运行超时"
			submit.Status = define.TimeLimitExceeded
		}
	}
	err = dao.DB.Create(submit).Error
	if err != nil {
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  "submit record create error:" + err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"data": gin.H{
			"data":   "submit success",
			"status": submit.Status,
			"msg":    msg,
		},
	})
}
