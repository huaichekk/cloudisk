package models

import (
	"gorm.io/gorm"
	"netdisk/commen/database"
)

type UserFile struct {
	gorm.Model
	UserID uint `gorm:"column:user_id;type:uint;not null;index;comment:用户ID"`
	User   User `gorm:"foreignKey:UserID;references:ID"` // 用户ID

	FileID *uint `gorm:"column:file_id;type:uint;index;comment:文件ID"`
	File   File  `gorm:"foreignKey:FileID;references:ID"` // 定义外键关系

	Type     int    `gorm:"column:type;type:tinyint;not null;default:1;comment:类型(1-文件 2-文件夹)"`
	FatherID int    `gorm:"column:father_id;type:int;not null;default:-1;index;comment:父目录ID(-1表示根目录)"`
	Name     string `gorm:"column:name;type:varchar(255);comment:名称(文件夹直接使用此字段)"`
}

func (UserFile) TableName() string {
	return "user_files"
}

func (uf *UserFile) GetName() string {
	if uf.Type == 1 {
		return uf.File.Name
	} else {
		return uf.Name
	}
}

func (uf UserFile) GetSize() uint64 {
	if uf.Type == 1 {
		return uf.File.Size
	} else {
		return 0
	}
}

func init() {
	if err := database.DB().AutoMigrate(&UserFile{}); err != nil {
		panic(err)
	}
}
