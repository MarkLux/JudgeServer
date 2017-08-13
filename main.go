package main

import (
	"fmt"

	"github.com/MarkLux/JudgeServer/compiler"
	"github.com/MarkLux/JudgeServer/config"
)

func main() {
	// r := gin.Default()

	// r.GET("/ping", server.Ping)

	// r.Run(":8090")

	test()
}

func test() {
	str := compiler.Compile(config.CompileC.CompileConfig, "/home/src", "/home/output")
	fmt.Println(str)
}
