package server

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/MarkLux/JudgeServer/client"

	"github.com/MarkLux/JudgeServer/compiler"

	"github.com/MarkLux/JudgeServer/config"

	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

type UserInput struct {
	Src        string `json:"src" binding:"required"`
	Language   string `json:"language" binding:"required"`
	MaxCPUTime int    `json:"max_cpu_time" binding:"required"`
	MaxMemory  int    `json:"max_memory" binding:"required"`
	TestCaseID int    `json:"test_case_id" binding:"required"`
}

func Ping(c *gin.Context) {

	if !checkToken(c.GetHeader("Token")) {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"data": "wrong token",
		})
		return
	}

	hostname, err := os.Hostname()

	if err != nil {
		hostname = ""
	}

	cpuPercent, _ := cpu.Percent(0, false)
	vmem, _ := mem.VirtualMemory()

	c.JSON(http.StatusOK, gin.H{
		"judger_version": "0.1.0",
		"hostname":       hostname,
		"cpu_core":       runtime.NumCPU(),
		"cpu":            cpuPercent,
		"memory":         vmem.UsedPercent,
		"action":         "pong",
	})
}

func Judge(c *gin.Context) {
	var input UserInput
	var err error

	if !checkToken(c.GetHeader("Token")) {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"data": "wrong token",
		})
		return
	}

	if err = c.BindJSON(&input); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": err.Error(),
		})
	}

	// generate submission id.

	submissionId := uuid.NewV4()

	// init submission dir

	var submissionDir string

	if submissionDir, err = initSubmissionEnv(submissionId.String()); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"data": err.Error(),
		})
		return
	}

	defer freeSubmissionEnv(submissionDir)

	var conf config.LanguageCompileConfig

	switch input.Language {
	case "c":
		conf = config.CompileC
		break
	case "cpp":
		conf = config.CompileCpp
		break
	case "java":
		conf = config.CompileJava
		break
	case "py2":
		conf = config.CompilePython2
		break
	default:
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"data": "invalid language support.",
		})
		return
	}

	// write source code into file.

	srcPath := filepath.Join(submissionDir, conf.SrcName)
	// log.Println("src path :", srcPath)
	var srcFile *os.File
	if srcFile, err = os.Create(srcPath); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"data": err.Error(),
		})
		return
	}
	defer srcFile.Close()

	_, err = srcFile.WriteString(input.Src)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"data": err.Error(),
		})
		return
	}

	// compile the source code

	var exePath string

	exePath, err = compiler.Compile(conf.CompileConfig, srcPath, submissionDir)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -2,
			"data": err.Error(),
		})
		return
	}

	// parse the run command
	// todo: parse the java's exe_dir and max_memory limit

	conf.RunConfig.Command.FillWith(map[string]string{
		"{exe_path}": exePath,
	})

	// run the client

	jc := client.JudgeClient{
		MaxCpuTime:    input.MaxCPUTime,
		MaxMemory:     input.MaxMemory,
		ExePath:       exePath,
		TestCaseId:    input.TestCaseID,
		SubmissionDir: submissionDir,
		RunConf:       conf.RunConfig,
	}

	result, err := jc.Judge()

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -3,
			"data": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": result,
	})

}

func SyncSingle(c *gin.Context) {

}

func initSubmissionEnv(submissionID string) (string, error) {
	submissionDirPath := filepath.Join(config.SUBMISSION_DIR, submissionID)

	if err := os.Mkdir(submissionDirPath, 0777); err != nil {
		return "", err
	}

	return submissionDirPath, nil
}

func freeSubmissionEnv(submissionDirPath string) error {
	return os.RemoveAll(submissionDirPath)
}

func checkToken(token string) bool {
	localToken := os.Getenv("RPC_TOKEN")
	log.Println(localToken)
	log.Println(token)
	if token != localToken {
		return false
	}
	return true
}
