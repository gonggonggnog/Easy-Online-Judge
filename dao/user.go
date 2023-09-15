package dao

import (
	"blog/models"
)

func GetPasswd(username string) (models.UserBasic, error) {
	var data models.UserBasic
	err := DB.Model(new(models.UserBasic)).Where("name=?", username).First(&data).Error // 查询
	return data, err
}
