package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/MarkLux/JudgeServer/client"

	"github.com/MarkLux/JudgeServer/compiler"
	"github.com/MarkLux/JudgeServer/config"
)

const (
	SUBMISSION_ROOT = "/home/judge/submissions"
	TESTCASE_ROOT   = "/home/judge/testcases"
)

var (
	pid      int
	progname string
)

const (
	logFileName = "test.log"
)

var resultArray = struct {
	sync.RWMutex
	m map[string]client.JudgeResult
}{m: make(map[string]client.JudgeResult)}

type jRes struct {
	testcaseId string
	result     client.JudgeResult
}

func init() {
	pid = os.Getpid()
	paths := strings.Split(os.Args[0], "/")
	paths = strings.Split(paths[len(paths)-1], string(os.PathSeparator))
	progname = paths[len(paths)-1]

	runtime.MemProfileRate = 1
}

func saveHeapProfile() {
	runtime.GC()

	f, err := os.Create(fmt.Sprintf("heap_%s_%d_%s.prof", progname, pid, time.Now().Format("2006_01_02_03_04_05")))
	if err != nil {
		return
	}
	defer f.Close()
	pprof.Lookup("heap").WriteTo(f, 1)
}

func Test_Judge(t *testing.T) {

	// defer saveHeapProfile()

	testcasesDirs, err := getDirs(TESTCASE_ROOT)

	if err != nil {
		t.Error("load testcases error!")
	}

	jobCh := make(chan string, 1000)
	resCh := make(chan jRes, 1000)

	for w := 0; w < 1; w++ {
		go worker(w, jobCh, resCh)
	}

	for _, dir := range testcasesDirs {
		jobCh <- dir
	}

	close(jobCh)

	for i := 0; i < len(testcasesDirs); i++ {
		rs := <-resCh

		if len(rs.result.UnPassed) > 0 {
			resultArray.Lock()
			resultArray.m[rs.testcaseId] = rs.result
			resultArray.Unlock()
			// printJRes()
			// close(resCh)
			// log.Println("ALL THINGS STTOPED!")
			// t.FailNow()
			// return
		} else {
			resultArray.Lock()
			resultArray.m[rs.testcaseId] = rs.result
			resultArray.Unlock()
		}
	}

	// resultArray.RLock()
	// for k, v := range resultArray.m {
	// 	log.Println(k, " : ", v)
	// }
	// resultArray.RUnlock()

	printJRes()

	return

}

func printJRes() {
	for k, v := range resultArray.m {
		resultArray.RLock()
		if len(v.Passed) <= 0 {
			log.Printf("%s test failed! result:\n%#v", k, v)
		} else {
			log.Printf("%s test ok!\n", k)
		}
		resultArray.RUnlock()
	}
}

func worker(id int, inCh <-chan string, outCh chan<- jRes) {
	for t := range inCh {
		log.Println("routine ", id, " processing testcase ", t)
		outCh <- jRes{
			testcaseId: t,
			result:     judgeTestcase(t),
		}
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

func judgeTestcase(testcaseId string) client.JudgeResult {

	var conf config.LanguageCompileConfig
	conf = config.CompileCpp

	log.Println("start test problem ", testcaseId)
	// 初始化变量
	submissionPath := filepath.Join(SUBMISSION_ROOT, testcaseId)
	srcFilePath := filepath.Join(submissionPath, "ac", "ac.cpp")
	// 编译
	log.Println("start compiling ...")

	exePath, err := compiler.Compile(conf.CompileConfig, srcFilePath, submissionPath)

	if err != nil {
		log.Println("compile error: ", err.Error())
		return client.JudgeResult{}
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
		return client.JudgeResult{}
	}

	// fmt.Printf("judge result of %s:\n%#v\n", testcaseId, result)

	logFile, er := os.Create(testcaseId + ".log")
	defer logFile.Close()
	if er != nil {
		log.Fatalln("fail to create log file")
	}

	logFile.WriteString(fmt.Sprintf("judge result of %s:\n%#v\n", testcaseId, result))

	return result
}
