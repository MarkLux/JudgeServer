package main

import (
	"github.com/MarkLux/JudgeServer/server"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/ping", server.Ping)

	r.POST("/judge", server.Judge)

	r.Run(":8090")
}
