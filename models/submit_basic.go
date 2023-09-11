package models

import "gorm.io/gorm"

type SubmitBasic struct {
	gorm.Model
	Identity        string         `gorm:"column:identity;type:varchar(36);" json:"identity"`
	ProblemIdentity string         `gorm:"column:problem_identity;type:varchar(36);" json:"problem_identity"`
	UserIdentity    string         `gorm:"column:user_identity;type:varchar(36);" json:"user_identity"`
	Status          int            `gorm:"column:status;type:tinyint;" json:"status"` //状态 -为等待判断，1为答案正确，2为答案错误，3为超时，4为超内存
	Path            string         `gorm:"column:path;type:varchar(255);" json:"path"`
	ProblemBasic    *ProblemsBasic `gorm:"foreignKey:identity;references:problem_identity;" json:"problem_basic"`
	UserBasic       *UserBasic     `gorm:"foreignKey:identity;references:user_identity;" json:"user_basic"`
}

func (table *SubmitBasic) TableName() string {
	return "submit_basic"
}
