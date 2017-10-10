package server

import (
	"log"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/MarkLux/JudgeServer/client"

	"github.com/MarkLux/JudgeServer/compiler"
	"github.com/MarkLux/JudgeServer/config"
)

const (
	SUBMISSION_ROOT = "/home/judge/submissions"
	TESTCASE_ROOT   = "/home/judge/testcases"
)

func Test_Judge(t *testing.T) {
	judgeTestcase("1002")
}

func judgeTestcase(testcaseId string) {

	var conf config.LanguageCompileConfig
	conf = config.CompileCpp

	log.Println("start test problem ", testcaseId)
	// 初始化变量
	submissionPath := filepath.Join(SUBMISSION_ROOT, testcaseId)
	srcFilePath := filepath.Join(submissionPath, "ac.cpp")
	// 编译
	log.Println("start compiling ...")

	exePath, err := compiler.Compile(conf.CompileConfig, srcFilePath, submissionPath)

	if err != nil {
		log.Println("compile error: ", err.Error())
		return
	}

	tId, _ := strconv.Atoi(testcaseId)

	conf.RunConfig.Command.FillWith(map[string]string{
		"{exe_path}": exePath,
		"{exe_dir}":  submissionPath,
	})

	jc := client.JudgeClient{
		MaxCpuTime:    5000,
		MaxMemory:     256 * 1024 * 1024,
		ExePath:       exePath,
		TestCaseId:    tId,
		SubmissionDir: submissionPath,
		RunConf:       config.CompileCpp.RunConfig,
	}

	result, err := jc.Judge()

	if err != nil {
		log.Println("run time error: " + err.Error())
		return
	}

	log.Printf("judge result:\n%#v\n", result)
	return
}
