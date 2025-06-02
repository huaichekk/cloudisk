package alive

import (
	"github.com/gin-gonic/gin"
	"netdisk/utils"
)

func AliveHandle() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		utils.OkWithMessage(ctx, "pong")
	}
}
