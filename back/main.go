package main

import (
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	r := gin.Default()
	//注册中间件

	//注册路由
	RegisterRouter(r)

	log.Fatalln(r.Run(":8080"))
}
