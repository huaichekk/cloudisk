package main

import (
	"github.com/gin-gonic/gin"
	"netdisk/controller/ai"
	"netdisk/controller/alive"
	"netdisk/controller/file"
	"netdisk/controller/user"
	"netdisk/middleware"
)

func RegisterRouter(r *gin.Engine) {
	api := r.Group("/api")
	api.Use(middleware.CorsMiddleWare())
	//无权限路由
	api.GET("/ping", alive.AliveHandle())
	userGroup := api.Group("/user")
	{
		userGroup.POST("/login", user.LoginHandle())
		userGroup.POST("/register", user.RegisterHandle())
	}

	//有权限路由
	authApi := api.Group("")
	authApi.Use(middleware.AuthMiddle())
	{
		authApi.GET("/auth_ping", alive.AuthAliveHandle())
		fileGroup := authApi.Group("/file")
		{
			fileGroup.GET("", file.GetFilesHandle())
			fileGroup.POST("/init", file.UploadInitHandle())
			fileGroup.POST("/chunk", file.UploadChunkHandle())
			fileGroup.POST("/merge", file.MergeHandle())
			fileGroup.GET("/download", file.DownLoadHandle())
			fileGroup.POST("/mkdir", file.MakeDir())
			fileGroup.POST("/remove", file.Remove())
		}
		aiGroup := authApi.Group("/ai")
		{
			aiGroup.POST("/chat", ai.ChatAI())
		}
	}
}
