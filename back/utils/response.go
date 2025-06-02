package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func OkWithData(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    data,
		"message": "success",
	})
}
func OkWithMessage(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    nil,
		"message": msg,
	})
}

func ServerFailWithMessage(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusOK, gin.H{
		"code":    500,
		"message": msg,
	})
}
func ClientFailWithMessage(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusOK, gin.H{
		"code":    400,
		"message": msg,
	})
}

func FailWithData(ctx *gin.Context, code int, data interface{}, message string) {
	ctx.JSON(code, gin.H{
		"code":    code,
		"data":    data,
		"message": message,
	})
}
