package file

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"netdisk/commen/database"
	"netdisk/models"
	"netdisk/utils"
	"os"
	"path/filepath"
)

func Remove() gin.HandlerFunc {
	//Post http:127.0.0.1:8080/api/files/remove
	return func(c *gin.Context) {
		var req struct {
			Id   int `json:"user_file_id"`
			Type int `json:"type"`
		}
		if err := c.BindJSON(&req); err != nil {
			utils.ClientFailWithMessage(c, err.Error())
			return
		}
		if req.Type == 1 { //文件
			var uf models.UserFile
			if err := database.DB().Where("id=?", req.Id).
				Preload("File").
				First(&uf).Error; err != nil {
				utils.ServerFailWithMessage(c, err.Error())
				return
			}
			if err := database.DB().Unscoped().Where("id=?", req.Id).Delete(&models.UserFile{}).Error; err != nil {
				utils.ServerFailWithMessage(c, err.Error())
				return
			}
			if err := database.DB().Unscoped().Where("id=?", uf.File.ID).Delete(&models.File{}).Error; err != nil {
				utils.ServerFailWithMessage(c, err.Error())
				return
			}
			if err := removeFile(uf.File.Hash); err != nil {
				utils.ServerFailWithMessage(c, err.Error())
				return
			}

		} else if req.Type == 2 { //文件夹
			var uf models.UserFile
			err := database.DB().Where("father_id=?", req.Id).First(&uf).Error
			if err == nil {
				utils.ClientFailWithMessage(c, "该文件夹下在存在文件,不可删除")
				return
			}
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				utils.ServerFailWithMessage(c, err.Error())
				return
			}
			if err := database.DB().Unscoped().Where("id=?", req.Id).Delete(&models.UserFile{}).Error; err != nil {
				utils.ServerFailWithMessage(c, err.Error())
				return
			}
		}
		utils.OkWithMessage(c, "删除成功")
		return
	}
}

func removeFile(fileHash string) error {
	path := filepath.Join(UPLOAD_PATH, fileHash)
	if err := os.Remove(path); err != nil {
		return err
	}
	return nil
}
