package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
	"time"
)

var (
	once sync.Once
	db   *gorm.DB
)

func DB() *gorm.DB {
	if db == nil {
		once.Do(func() {
			var err error
			dsn := "root:320930@tcp(124.70.56.84:3306)/netdisk?charset=utf8mb4&parseTime=True&loc=Local"
			db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
			if err != nil {
				panic(err)
			}
			sqlDB, _ := db.DB()
			sqlDB.SetMaxIdleConns(10)  // 空闲连接数
			sqlDB.SetMaxOpenConns(100) // 最大打开连接数
			sqlDB.SetConnMaxLifetime(time.Hour)
		})
	}
	return db
}
