package dao

import (
	"blog/models"
	"gorm.io/gorm"
)

func GetProblemList(keyword, categoryIdentity string) *gorm.DB {
	tx := DB.Model(new(models.ProblemsBasic)).Preload("ProblemCategories").
		Preload("ProblemCategories.CategoryBasic").
		Where("title like ? OR content like ?", "%"+keyword+"%", "%"+keyword+"%") // 模糊查询
	if categoryIdentity != "" {
		tx.Joins("RIGHT JOIN problem_category pc on pc.problem_id = problems_basic.id").
			Where("pc.category_id = any(SELECT cb.id FROM category_basic cb WHERE cb.identity = ? )",
				categoryIdentity) // 根据分类标识查询
	}
	return tx
}
