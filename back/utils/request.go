package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func GetUserID(ctx *gin.Context) (*uint, error) {
	if value, ok := ctx.Get("userID"); !ok {
		return nil, fmt.Errorf("用户名ID获取失败")
	} else {
		return value.(*uint), nil
	}
}
