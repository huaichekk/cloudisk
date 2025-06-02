package file

import (
	"github.com/gin-gonic/gin"
	"netdisk/commen/database"
	"netdisk/models"
	"netdisk/utils"
)

func GetFilesHandle() gin.HandlerFunc {
	//Get http:127.0.0.1:8080/api/files?father_id=-1
	return func(context *gin.Context) {
		father_id := context.Query("father_id")
		id, err := utils.GetUserID(context)
		if err != nil {
			utils.ClientFailWithMessage(context, err.Error())
			return
		}
		var userFiles []models.UserFile
		if err := database.DB().
			Where("father_id = ?", father_id).
			Where("user_id=?", id).
			Preload("File").
			Find(&userFiles).
			Error; err != nil {
			utils.ServerFailWithMessage(context, err.Error())
			return
		}
		respData := make([]map[string]interface{}, 0)
		for _, uf := range userFiles {
			item := map[string]interface{}{
				"user_file_id": uf.ID,
				"file_name":    uf.GetName(),
				"file_size":    uf.GetSize(),
				"update_at":    uf.UpdatedAt,
				"type":         uf.Type,
			}
			respData = append(respData, item)
		}
		utils.OkWithData(context, respData)
		return
	}
}
