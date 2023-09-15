package dao

import (
	"blog/models"
	"gorm.io/gorm"
)

func GetCategoryList(keyword string) *gorm.DB {
	return DB.Model(new(models.CategoryBasic)).Where("name like ?", "%"+keyword+"%") // 模糊查询
}
