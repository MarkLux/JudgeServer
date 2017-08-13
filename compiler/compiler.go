package compiler

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/MarkLux/JudgeServer/config"
	"github.com/MarkLux/Judger_GO"
)

/*
* compile the program in the judger sandbox.
* return the executable file path
 */

func Compile(compileConfig config.CompileConfig, srcPath string, outputDir string) string {
	exePath, _ := filepath.Abs(filepath.Join(outputDir, compileConfig.ExeName))
	// replace param then build the real compile command.
	replacements := map[string]string{
		"{src_path}": srcPath,
		"{exe_dir}":  outputDir,
		"{exe_path}": exePath,
	}
	command := compileConfig.CompileCommand.FillWith(replacements)
	// file compiler.out contains the ouput of compile progarm's output (rewrite).
	compilerOut := filepath.Join(outputDir, "compiler.out")

	// split the command into execute path and args.
	spilts := strings.Split(command, " ")
	//parse args
	var args [judger.ARGS_MAX_NUMBER]string
	for i, split := range spilts[1:] {
		args[i] = split
	}
	//parse envs
	var envs [judger.ENV_MAX_NUMBER]string
	envs[0] = "PATH=" + os.Getenv("PATH")

	result := judger.JudgerRun(judger.Config{
		MaxCpuTime:       compileConfig.MaxCpuTime,
		MaxRealTime:      compileConfig.MaxRealTime,
		MaxMemory:        compileConfig.MaxMemory,
		MaxStack:         128 * 1024 * 1024,
		MaxProcessNumber: judger.UNLIMITED,
		ExePath:          spilts[0],
		InputPath:        srcPath,
		OutputPath:       compilerOut,
		ErrorPath:        compilerOut,
		Args:             args,
		Env:              envs,
		LogPath:          config.COMPILER_LOG_PATH,
		SecCompRuleName:  "",
		Uid:              config.COMPILER_USER_UID,
		Gid:              config.COMPILER_GROUP_UID,
	})

	return string(result.Result)
}
