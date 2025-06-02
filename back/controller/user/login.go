package user

import (
	"github.com/gin-gonic/gin"
	"netdisk/commen/database"
	"netdisk/middleware"
	"netdisk/models"
	"netdisk/utils"
)

func LoginHandle() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		username := ctx.PostForm("username")
		password := ctx.PostForm("password")
		var user models.User
		if err := database.DB().Where("username=?", username).First(&user).Error; err != nil {
			utils.ServerFailWithMessage(ctx, err.Error())
			return
		}
		if user.Password != password {
			utils.ClientFailWithMessage(ctx, "密码错误，请重新输入")
			return
		}
		token, err := middleware.GetToken(user)
		if err != nil {
			utils.ServerFailWithMessage(ctx, err.Error())
			return
		}
		utils.OkWithData(ctx, map[string]string{
			"X-Token": token,
		})
		return
	}
}
