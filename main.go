package main

import (
	"github.com/MarkLux/JudgeServer/server"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/ping", server.Ping)

	r.Run(":8090")
}
