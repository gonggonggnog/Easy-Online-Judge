package models

import "gorm.io/gorm"

type ProblemsCategory struct {
	gorm.Model
	ProblemId     string         `gorm:"column:problem_id;type:varchar(36);" json:"problem_id"`
	CategoryId    string         `gorm:"column:category_id;type:varchar(36);" json:"category_id"`
	CategoryBasic *CategoryBasic `gorm:"foreignKey:id;references:category_id;" json:"category_basic"`
}

func (table *ProblemsCategory) TableName() string {
	return "problem_category"
}
