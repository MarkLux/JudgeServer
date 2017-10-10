package server

import (
	"io/ioutil"
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
	testcasesDirs, err := getDirs(TESTCASE_ROOT)

	if err != nil {
		t.Error("load testcases error!")
	}

	jobCh := make(chan string, 100)
	resCh := make(chan bool, 100)

	for w := 0; w < 10; w++ {
		go worker(w, jobCh, resCh)
	}

	for _, dir := range testcasesDirs {
		jobCh <- dir
	}

	close(jobCh)

	resultArray := make(map[string]string)

	for i := 0; i < len(testcasesDirs); i++ {
		rs := <-resCh

		if rs {
			resultArray[testcasesDirs[i]] = "test ok!"
		} else {
			resultArray[testcasesDirs[i]] = "test failed!"
		}
	}

	for k, v := range resultArray {
		log.Println(k, " : ", v)
	}

}

func worker(id int, inCh <-chan string, outCh chan<- bool) {
	for t := range inCh {
		log.Println("routine ", id, " processing testcase ", t)
		outCh <- judgeTestcase(t)
	}
}

func getDirs(dirPth string) (files []string, err error) {
	files = make([]string, 0, 10)
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}

	for _, fi := range dir {
		if fi.IsDir() {
			files = append(files, fi.Name())
		}
	}
	return files, nil
}

func judgeTestcase(testcaseId string) bool {

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
		return false
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
		RunConf:       conf.RunConfig,
	}

	result, err := jc.Judge()

	if err != nil {
		log.Println("run time error: " + err.Error())
		return false
	}

	log.Printf("judge result:\n%#v\n", result)
	if len(result.UnPassed) == 0 {
		return true
	} else {
		return false
	}
}
