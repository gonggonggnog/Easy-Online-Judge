package models

import "gorm.io/gorm"

type ProblemsBasic struct {
	gorm.Model
	Identity          string              `gorm:"column:identity;type:varchar(36);" json:"identity"` //问题的标识
	ProblemCategories []*ProblemsCategory `gorm:"foreignKey:problem_id;references:id;"`              //分类id
	Title             string              `gorm:"column:title;type:varchar(255);" json:"title"`      //问题标题
	Content           string              `gorm:"column:content;type:text;" json:"content"`          //问题正文描述
	MaxRuntime        int                 `gorm:"column:max_runtime;type:int;" json:"max_runtime"`
	MaxMem            int                 `gorm:"column:max_mem;type:int;" json:"max_mem"`
}

func (table *ProblemsBasic) TableName() string {
	return "problems_basic"
}
