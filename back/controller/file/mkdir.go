package file

import (
	"github.com/gin-gonic/gin"
	"netdisk/commen/database"
	"netdisk/models"
	"netdisk/utils"
)

func MakeDir() gin.HandlerFunc {
	////post http:127.0.0.1:8080/api/files
	return func(c *gin.Context) {
		var req struct {
			FatherID int    `json:"father_id"`
			DirName  string `json:"dir_name"`
		}
		if err := c.BindJSON(&req); err != nil {
			utils.ClientFailWithMessage(c, err.Error())
			return
		}
		id, err := utils.GetUserID(c)
		if err != nil {
			utils.ClientFailWithMessage(c, err.Error())
			return
		}
		uf := &models.UserFile{
			UserID:   *id,
			FileID:   nil,
			Type:     2,
			FatherID: req.FatherID,
			Name:     req.DirName,
		}
		if err := database.DB().Create(&uf).Error; err != nil {
			utils.ServerFailWithMessage(c, err.Error())
			return
		}
		utils.OkWithMessage(c, "创建文件夹成功")
		return
	}
}
