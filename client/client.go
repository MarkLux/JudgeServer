package client

import (
	"bytes"
	"crypto/md5"
	"errors"
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
	CpuTime     int
	Result      judger.ResultCode
	RealTime    int
	Signal      int
	Error       judger.ErrorCode
	OutputMD5   string
	TestCaseNum int
}

type JudgeClient struct {
	RunConf       config.RunConfig
	ExePath       string
	MaxCpuTime    int
	MaxMemory     int
	SubmissionDir string
	TestCaseId    int
}

// func (*JudgeClient) Judge() (JudgeResult, error) {

// }

func (jc *JudgeClient) JudgeOne(testInPath string, testOutPath string) (userOutputMd5 string, res bool, err error) {
	commands := strings.Split(string(jc.RunConf.Command), " ")
	userOutputPath := filepath.Join(jc.SubmissionDir, "user.out")
	fmt.Println(userOutputPath)
	runResult := judger.JudgerRun(judger.Config{
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
	fmt.Printf("%#v", runResult)

	if runResult.Error != judger.SUCCESS {
		err = errors.New("Runtime Error, Code" + fmt.Sprintf("%#v", runResult))
		return
	}

	userOutputMd5, res = compareOutput(testOutPath, userOutputPath)
	err = nil
	return
}

func compareOutput(testOutPath string, userOutputPath string) (outputMd5 string, res bool) {
	testOut, _ := ioutil.ReadFile(testOutPath)
	// Linux files has a line break at the end defaultly,trim it!
	trimed := bytes.TrimRight(testOut, "\n")
	testOut, _ = ioutil.ReadFile(testOutPath)
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
