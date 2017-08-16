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

func (jc *JudgeClient) Judge() (judgeResult JudgeResult, err error) {
	// load testcases files
	testFiles, err := loadTestCases(jc.TestCaseId)
	if err != nil {
		return JudgeResult{}, err
	}

	// build channel
	caseCount := len(testFiles)
	caseCh := make(chan RunResult, caseCount)

	for _, t := range testFiles {
		testInPath := filepath.Join(config.TEST_CASE_DIR, strconv.Itoa(jc.TestCaseId), t+".in")
		testOutPath := filepath.Join(config.TEST_CASE_DIR, strconv.Itoa(jc.TestCaseId), t+".out")
		// create a new go routine for each testcase file
		go jc.judgeOne(caseCh, testInPath, testOutPath)
	}

	for i := 0; i < caseCount; i++ {
		rs := <-caseCh
		if rs.Result == judger.SUCCESS {
			judgeResult.Passed = append(judgeResult.Passed, rs)
		} else {
			judgeResult.UnPassed = append(judgeResult.UnPassed, rs)
		}
	}
	err = nil

	return
}

func (jc *JudgeClient) judgeOne(ch chan<- RunResult, testInPath string, testOutPath string) {
	commands := strings.Split(string(jc.RunConf.Command), " ")
	userOutputPath := filepath.Join(jc.SubmissionDir, "user.out")
	fmt.Println(testInPath)
	fmt.Println(testOutPath)
	fmt.Println(userOutputPath)
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
	ch <- runResult
	return
}

func compareOutput(testOutPath string, userOutputPath string) (outputMd5 string, res bool) {
	testOut, _ := ioutil.ReadFile(testOutPath)
	// Linux files has a line break at the end defaultly,trim it!
	trimed := bytes.TrimRight(testOut, "\n")
	testMD5 := md5.Sum(trimed)
	userOut, _ := ioutil.ReadFile(userOutputPath)
	userMD5 := md5.Sum(userOut)
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
