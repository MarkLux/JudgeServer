package client

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/MarkLux/JudgeServer/config"
	"github.com/MarkLux/JudgeServer/utils"
	"github.com/MarkLux/Judger_GO"
)

type JudgeResult struct {
	Passed   []RunResult
	UnPassed []RunResult
}

type RunResult struct {
	CpuTime   int
	Result    judger.ResultCode
	RealTime  int
	Memory    int
	Signal    int
	Error     judger.ErrorCode
	OutputMD5 string
}

type JudgeClient struct {
	RunConf       config.RunConfig
	ExePath       string
	MaxCpuTime    int
	MaxMemory     int
	SubmissionDir string
	TestCaseId    int
}

type TestCase struct {
	UserOutFilename string
	TestInPath      string
	TestOutPath     string
}

var instanceCount int

func (jc *JudgeClient) Judge() (judgeResult JudgeResult, err error) {

	// only for test
	// runtime.GOMAXPROCS(runtime.NumCPU())

	// load testcases files
	testFiles, err := loadTestCases(jc.TestCaseId)
	if err != nil {
		return JudgeResult{}, err
	}

	// build channel
	jobCh := make(chan TestCase, 100)
	resCh := make(chan RunResult, 100)

	// routine pool

	for w := 0; w < config.MAX_JUDGE_ROUTINES; w++ {
		go worker(jc, jobCh, resCh)
	}

	// add jobs

	for _, t := range testFiles {
		tc := TestCase{
			TestInPath:      filepath.Join(config.TEST_CASE_DIR, strconv.Itoa(jc.TestCaseId), t+".in"),
			TestOutPath:     filepath.Join(config.TEST_CASE_DIR, strconv.Itoa(jc.TestCaseId), t+".out"),
			UserOutFilename: t + ".out",
		}
		jobCh <- tc
	}

	close(jobCh)

	for i := 0; i < len(testFiles); i++ {
		rs := <-resCh
		if rs.Result == judger.SUCCESS {
			judgeResult.Passed = append(judgeResult.Passed, rs)
		} else {
			judgeResult.UnPassed = append(judgeResult.UnPassed, rs)
		}
	}

	err = nil

	return
}

func worker(jc *JudgeClient, inCh <-chan TestCase, outCh chan<- RunResult) {
	for t := range inCh {
		outCh <- jc.JudgeOne(t.TestInPath, t.TestOutPath, t.UserOutFilename)
	}
}

func (jc *JudgeClient) JudgeOne(testInPath string, testOutPath string, userOutFilename string) RunResult {
	commands := strings.Split(string(jc.RunConf.Command), " ")
	userOutputPath := filepath.Join(jc.SubmissionDir, userOutFilename)
	result := judger.JudgerRun(judger.Config{
		MaxCpuTime:       jc.MaxCpuTime,
		MaxMemory:        jc.MaxMemory,
		MaxStack:         128 * 1024 * 1024,
		MaxOutPutSize:    1024 * 1024 * 1024,
		MaxRealTime:      jc.MaxCpuTime * 3,
		MaxProcessNumber: judger.UNLIMITED,
		ExePath:          commands[0],
		InputPath:        testInPath,
		OutputPath:       userOutputPath,
		ErrorPath:        userOutputPath,
		Args:             commands,
		Env:              append(jc.RunConf.Env, "PATH="+os.Getenv("PATH")),
		LogPath:          config.JUDGE_RUN_LOG_PATH,
		SecCompRuleName:  jc.RunConf.SeccompRule,
		Uid:              config.RUN_USER_UID,
		Gid:              config.RUN_GROUP_UID,
	})

	// if result.Error != judger.SUCCESS {
	// 	err = errors.New("Runtime Error, Code" + fmt.Sprintf("%#v", result))
	// 	return
	// }

	userOutputMd5, res := compareOutput(testOutPath, userOutputPath)
	if !res {
		// set to wrong answer
		result.Result = -1
	}

	runResult := RunResult{
		CpuTime:   result.CpuTime,
		RealTime:  result.RealTime,
		Memory:    result.Memory,
		Result:    result.Result,
		Error:     result.Error,
		OutputMD5: userOutputMd5,
		Signal:    result.Signal,
	}

	return runResult
}

func compareOutput(testOutPath string, userOutputPath string) (outputMd5 string, res bool) {
	testOut, _ := ioutil.ReadFile(testOutPath)
	// Linux files has a line break at the end defaultly,trim it!
	trimed := bytes.TrimRight(testOut, "\n")
	testMD5 := md5.Sum(trimed)
	userOut, _ := ioutil.ReadFile(userOutputPath)
	trimed = bytes.TrimRight(userOut, "\n")
	userMD5 := md5.Sum(trimed)
	outputMd5 = fmt.Sprintf("%x", userMD5)
	res = bool(testMD5 == userMD5)
	return
}

func loadTestCases(testCaseID int) (fileName []string, err error) {
	path := filepath.Join(config.TEST_CASE_DIR, strconv.Itoa(testCaseID))
	inFileNames, err := utils.GetFileNames(path, ".in")
	if err != nil {
		return []string{}, err
	}
	for _, fiName := range inFileNames {
		if _, e := os.Stat(filepath.Join(path, fiName+".out")); e == nil {
			fileName = append(fileName, fiName)
		}
	}
	return fileName, nil
}
