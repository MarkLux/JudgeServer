package main

import (
	"fmt"

	"github.com/MarkLux/JudgeServer/client"
	"github.com/MarkLux/JudgeServer/config"
	"github.com/MarkLux/JudgeServer/server"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/ping", server.Ping)

	r.POST("/judge", server.Judge)

	r.Run(":8070")

	// test()
}

func test() {

	conf := config.CompileC.RunConfig

	conf.Command.FillWith(map[string]string{
		"{exe_path}": "/home/judge/output/main",
	})

	jc := client.JudgeClient{
		MaxCpuTime:    1000,
		MaxMemory:     128 * 1024 * 1024,
		RunConf:       conf,
		TestCaseId:    1001,
		SubmissionDir: "/home/judge/user/1001",
	}

	rs, _ := jc.Judge()

	for _, p := range rs.Passed {
		fmt.Println("passed: ")
		fmt.Printf("%#v", p)
	}
	for _, u := range rs.UnPassed {
		fmt.Println("unpassed: ")
		fmt.Printf("%#v", u)
	}
}
