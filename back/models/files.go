package models

import (
	"gorm.io/gorm"
	"netdisk/commen/database"
)

type File struct {
	gorm.Model

	// 文件名
	Name string `gorm:"column:name;type:varchar(255);not null;index;comment:文件名"`

	// 文件大小（字节）
	Size uint64 `gorm:"column:size;type:bigint unsigned;not null;default:0;comment:文件大小(字节)"`

	// 文件哈希值（用于去重校验）
	Hash string `gorm:"column:hash;type:varchar(64);not null;uniqueIndex;comment:文件哈希值(MD5)"`
}

func (f *File) TableName() string {
	return "files"
}
func init() {
	if err := database.DB().AutoMigrate(&File{}); err != nil {
		panic(err)
	}
}
