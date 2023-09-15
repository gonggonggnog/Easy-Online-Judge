package models

import "gorm.io/gorm"

type TestCase struct {
	gorm.Model
	Identity        string `gorm:"column:identity;type:varchar(36);" json:"identity"`                 // 测试样例的标识
	ProblemIdentity string `gorm:"column:problem_identity;type:varchar(36);" json:"problem_identity"` // 问题的标识
	Input           string `gorm:"column:input;type:text;" json:"input"`                              // 测试样例的输入
	Output          string `gorm:"column:output;type:text;" json:"output"`                            // 测试样例的输出
}

func (table *TestCase) TableName() string {
	return "test_case"
}
