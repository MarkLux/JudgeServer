package config

import (
	"strings"

	"github.com/MarkLux/Judger_GO"
)

type CommandStr string

func (c *CommandStr) FillWith(param map[string]string) string {
	for k, v := range param {
		*c = CommandStr(strings.Replace(string(*c), k, v, -1))
	}
	return string(*c)
}

type CompileConfig struct {
	SrcName string // source code file name
	ExeName string // executable file name (after compile)
	// compile will run in the judger sandbox with configs below.
	MaxCpuTime  int
	MaxRealTime int
	MaxMemory   int
	// run the command to compile the code.
	CompileCommand CommandStr
}

type RunConfig struct {
	Command     CommandStr
	SeccompRule string
	Env         []string
}

type LanguageCompileConfig struct {
	CompileConfig
	RunConfig
}

var DefaultEnv = []string{
	"LANG=en_US.UTF-8",
	"LANGUAGE=en_US:en",
	"LC_ALL=en_US.UTF-8",
}

var CompileC = LanguageCompileConfig{
	CompileConfig: CompileConfig{
		SrcName:        "main.c",
		ExeName:        "main",
		MaxCpuTime:     3000,
		MaxRealTime:    5000,
		MaxMemory:      judger.UNLIMITED,
		CompileCommand: "/usr/bin/gcc -DONLINE_JUDGE -O2 -w -fmax-errors=3 -std=c99 {src_path} -lm -o {exe_path}",
	},
	RunConfig: RunConfig{
		Command:     "{exe_path}",
		SeccompRule: "c_cpp",
		Env:         DefaultEnv,
	},
}

var CompileCpp = LanguageCompileConfig{
	CompileConfig: CompileConfig{
		SrcName:        "main.cpp",
		ExeName:        "main",
		MaxCpuTime:     3000,
		MaxRealTime:    5000,
		MaxMemory:      judger.UNLIMITED,
		CompileCommand: "/usr/bin/g++ -DONLINE_JUDGE -O2 -w -fmax-errors=3 -std=c++11 {src_path} -lm -o {exe_path}",
	},
	RunConfig: RunConfig{
		Command:     "{exe_path}",
		SeccompRule: "c_cpp",
		Env:         DefaultEnv,
	},
}

var CompileJava = LanguageCompileConfig{
	CompileConfig: CompileConfig{
		SrcName:        "Main.java",
		ExeName:        "Main",
		MaxCpuTime:     3000,
		MaxRealTime:    5000,
		MaxMemory:      judger.UNLIMITED,
		CompileCommand: "/usr/bin/javac {src_path} -d {exe_dir} -encoding UTF8",
	},
	RunConfig: RunConfig{
		Command:     "/usr/bin/java -cp {exe_dir} -Xss1M -XX:MaxPermSize=16M -XX:PermSize=8M -Xms16M -Xmx{max_memory}k -Djava.security.manager -Dfile.encoding=UTF-8 -Djava.security.policy==/etc/java_policy -Djava.awt.headless=true Main",
		SeccompRule: "none",
		Env:         append(DefaultEnv, "MALLOC_ARENA_MAX=1"),
	},
}

var CompilePython2 = LanguageCompileConfig{
	CompileConfig: CompileConfig{
		SrcName:        "solution.py",
		ExeName:        "solution.pyc",
		MaxCpuTime:     3000,
		MaxRealTime:    5000,
		MaxMemory:      judger.UNLIMITED,
		CompileCommand: "/usr/bin/python -m py_compile {src_path}",
	},
	RunConfig: RunConfig{
		Command:     "/usr/bin/python {exe_path}",
		SeccompRule: "none",
		Env:         DefaultEnv,
	},
}
