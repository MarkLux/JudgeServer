package compiler

import (
	"fmt"
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
	exePath := filepath.Join(outputDir, compileConfig.ExeName)
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
	args := strings.Split(command, " ")

	//parse args

	result := judger.JudgerRun(judger.Config{
		MaxCpuTime:       compileConfig.MaxCpuTime,
		MaxRealTime:      compileConfig.MaxRealTime,
		MaxMemory:        compileConfig.MaxMemory,
		MaxStack:         128 * 1024 * 1024,
		MaxOutPutSize:    judger.UNLIMITED,
		MaxProcessNumber: judger.UNLIMITED,
		ExePath:          args[0],
		InputPath:        srcPath,
		OutputPath:       compilerOut,
		ErrorPath:        compilerOut,
		Args:             args,
		Env:              []string{"PATH=" + os.Getenv("PATH")},
		LogPath:          config.COMPILER_LOG_PATH,
		SecCompRuleName:  "none",
		Uid:              config.COMPILER_USER_UID,
		Gid:              config.COMPILER_GROUP_UID,
	})

	// debug output

	fmt.Printf("%#v", result)

	return string(result.Result)
}
