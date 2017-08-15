package main

import (
	"fmt"

	"github.com/MarkLux/JudgeServer/config"

	"github.com/MarkLux/JudgeServer/client"
)

func main() {
	// r := gin.Default()

	// r.GET("/ping", server.Ping)

	// r.Run(":8090")

	test()
}

func test() {

	conf := config.CompileC.RunConfig

	conf.Command.FillWith(map[string]string{
		"{exe_path}": "/home/judge/output/main",
	})

	fmt.Println(conf.Command)

	jc := client.JudgeClient{
		MaxCpuTime:    1000,
		MaxMemory:     128 * 1024 * 1024,
		RunConf:       conf,
		TestCaseId:    1001,
		SubmissionDir: "/home/judge/user/1001",
	}

	m5, res, err := jc.JudgeOne("/home/judge/testcases/1001/1.in", "/home.judge/testcase/1001/1.out")

	println("m5 = ", m5, "; res = ", res, " ;err = ", err)
}
