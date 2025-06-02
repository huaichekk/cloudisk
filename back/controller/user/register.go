package user

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"netdisk/commen/database"
	"netdisk/models"
	"netdisk/utils"
)

func RegisterHandle() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		username := ctx.PostForm("username")
		password := ctx.PostForm("password")
		var user models.User
		err := database.DB().Select("password").Where("username=?", username).First(&user).Error
		if err == nil { //用户名重复
			utils.ClientFailWithMessage(ctx, "用户名重复，请重新输入")
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) { //其他错误
			utils.ServerFailWithMessage(ctx, err.Error())
		} else { //没有查找到用户名相同的记录
			u := models.User{
				UserName: username,
				Password: password,
			}
			if err = database.DB().Create(&u).Error; err != nil {
				utils.ServerFailWithMessage(ctx, err.Error())
			} else {
				utils.OkWithMessage(ctx, "注册成功")
			}
		}
	}
}
