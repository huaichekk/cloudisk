package models

import (
	"gorm.io/gorm"
	"netdisk/commen/database"
)

type User struct {
	gorm.Model
	UserName string `gorm:"uniqueIndex;type:varchar(255);not_null;column:username"`
	Password string `gorm:"type:varchar(255);not_null;column:password"`
	IsVIP    bool   `gorm:"default:false"`
}

func (User) TableName() string {
	return "users"
}

func init() {
	if err := database.DB().AutoMigrate(&User{}); err != nil {
		panic(err)
	}
}
