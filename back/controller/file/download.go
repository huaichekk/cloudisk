package file

import (
	"github.com/gin-gonic/gin"
	"netdisk/commen/database"
	"netdisk/models"
	"netdisk/utils"
	"path/filepath"
)

func DownLoadHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Query("user_file_id")
		var uf models.UserFile
		if err := database.DB().Where("id=?", id).
			Preload("File").
			First(&uf).Error; err != nil {
			utils.ServerFailWithMessage(c, err.Error())
			return
		}
		filePath := filepath.Join(UPLOAD_PATH, uf.File.Hash)
		c.File(filePath)
		return
	}
}
