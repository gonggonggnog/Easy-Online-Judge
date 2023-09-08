package dao

import (
	"blog/models"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB = Init()

func Init() *gorm.DB {
	dsn := "root:123456@tcp(47.115.224.170:3306)/gin_oj?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("连接失败")
		fmt.Println(err)

	}
	return db
}

func GetProblemList(keyword, categoryIdentity string) *gorm.DB {
	tx := DB.Model(new(models.ProblemsBasic)).Preload("ProblemCategories").Preload("ProblemCategories.CategoryBasic").Where("title like ? OR content like ?", "%"+keyword+"%", "%"+keyword+"%")
	if categoryIdentity != "" {
		tx.Joins("RIGHT JOIN problem_category pc on pc.problem_id = problems_basic.id").
			Where("pc.category_id = any(SELECT cb.id FROM category_basic cb WHERE cb.identity = ? )", categoryIdentity)
	}
	return tx
}
