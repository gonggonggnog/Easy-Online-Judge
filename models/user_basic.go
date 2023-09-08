package models

import "gorm.io/gorm"

type UserBasic struct {
	gorm.Model
	Identity string `gorm:"column:identity;type:varchar(36);" json:"identity"`
	Name     string `gorm:"column:name;type:varchar(100);" json:"name"`
	Password string `gorm:"column:password;type:varchar(100);" json:"password"`
	Email    string `gorm:"column:email;type:varchar(36);" json:"email"`
	Phone    string `gorm:"column:phone;type:varchar(20);" json:"phone"`
}

func (table *UserBasic) TableName() string {
	return "user_basic"
}
