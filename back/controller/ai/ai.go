package ai

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"netdisk/utils"
)

func ChatAI() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		msg := ctx.PostForm("message")
		fmt.Println("接收到的信息：", msg)
		id, err := utils.GetUserID(ctx)
		if err != nil {
			utils.ClientFailWithMessage(ctx, err.Error())
			return
		}
		message, err := utils.ChatMessage(msg, int(*id))
		if err != nil {
			utils.ServerFailWithMessage(ctx, err.Error())
			return
		}
		log.Println(message)
		utils.OkWithData(ctx, map[string]interface{}{
			"message": message,
		})
	}
}
