package dao

import (
	"blog/models"
	"gorm.io/gorm"
)

func GetSubmitList(userIdentity, problemIdentity string, status int) *gorm.DB {
	tx := DB.Model(new(models.SubmitBasic)).Preload("ProblemBasic").Preload("UserBasic")
	if userIdentity != "" {
		tx.Where("user_identity=?", userIdentity)
	}
	if problemIdentity != "" {
		tx.Where("problem_identity=?", problemIdentity)
	}
	if status != 0 {
		tx.Where("status=?", status)
	}
	return tx
}
