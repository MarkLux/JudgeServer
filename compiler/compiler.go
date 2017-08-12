package compiler

import (
	"path/filepath"

	"github.com/MarkLux/JudgeServer/config"
)

/*
* compile the program in the judger sandbox.
* return the executable file path
 */

func Compile(compileConfig config.CompileConfig, srcPath string, outputDir string) string {
	exePath, _ := filepath.Abs(filepath.Join(outputDir, compileConfig.ExeName))
	command := compileConfig.CompileCommand
}
