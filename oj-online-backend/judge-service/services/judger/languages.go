package judger

type CompileConfig struct {
	SrcName        string `json:"src_name"`
	ExeName        string `json:"exe_name"`
	MaxCpuTime     int64  `json:"max_cpu_time"`
	MaxRealTime    int64  `json:"max_real_time"`
	MaxMemory      int64  `json:"max_memory"`
	CompileCommand string `json:"compile_command"`
}

type RunConfig struct {
	Command              string   `json:"command"`
	SeccompRule          string   `json:"seccomp_rule"`
	Env                  []string `json:"env"`
	MemoryLimitCheckOnly int      `json:"memory_limit_check_only"`
}

type LangConfig struct {
	CompileConfig CompileConfig `json:"compile"`
	RunConfig     RunConfig     `json:"run"`
}

var DefaultEnv = []string{"LANG=en_US.UTF-8", "LANGUAGE=en_US:en", "LC_ALL=en_US.UTF-8"}

var CPPLangConfig = &LangConfig{
	CompileConfig: CompileConfig{
		SrcName:        "main.cpp",
		ExeName:        "main",
		MaxCpuTime:     3000,
		MaxRealTime:    5000,
		MaxMemory:      128 * 1024 * 1024,
		CompileCommand: "/usr/bin/g++ -DONLINE_JUDGE -O2 -w -fmax-errors=3 -std=c++11 {src_path} -lm -o {exe_path}",
	},
	RunConfig: RunConfig{
		Command:     "{exe_path}",
		SeccompRule: "c_cpp",
		Env:         DefaultEnv,
	},
}
