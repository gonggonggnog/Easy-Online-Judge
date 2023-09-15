package models

import "gorm.io/gorm"

type ProblemsBasic struct {
	gorm.Model
	Identity          string              `gorm:"column:identity;type:varchar(36);" json:"identity"`             //问题的标识
	ProblemCategories []*ProblemsCategory `gorm:"foreignKey:problem_id;references:id" json:"problem_categories"` //分类id
	Title             string              `gorm:"column:title;type:varchar(255);" json:"title"`                  //问题标题
	Content           string              `gorm:"column:content;type:text;" json:"content"`                      //问题正文描述
	MaxRuntime        int                 `gorm:"column:max_runtime;type:int;" json:"max_runtime"`
	MaxMem            int                 `gorm:"column:max_mem;type:int;" json:"max_mem"`
	PassNum           int64               `gorm:"column:pass_num;type:int(11);" json:"pass_num"`     // 通过的次数
	SubmitNum         int64               `gorm:"column:submit_num;type:int(11);" json:"submit_num"` // 提交次数
	TestCase          []*TestCase         `gorm:"foreignKey:problem_identity;references:identity" json:"test_case"`
}

func (table *ProblemsBasic) TableName() string {
	return "problems_basic"
}
