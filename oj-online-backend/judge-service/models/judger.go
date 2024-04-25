package models

type CompileConfig struct {
	SrcName     int64  `json:"src_name"`
	ExeName     int64  `json:"exe_name"`
	MaxCpuTime  int64  `json:"max_cpu_time"`
	MaxRealTime int64  `json:"max_real_time"`
	MaxMemory   int64  `json:"max_memory"`
	CompileExe  string `json:"compiler_exe"`
	CompileArgs string `json:"compile_args"`
}

type RunConfig struct {
	SeccompRuleName      string   `json:"seccomp_rule_name"`
	Env                  []string `json:"env"`
	MemoryLimitCheckOnly int      `json:"memory_limit_check_only"`
	MaxCpuTime           int64    `json:"max_cpu_time"`
	MaxRealTime          int64    `json:"max_real_time"`
	MaxMemory            int64    `json:"max_memory"`
}

type LangConfig struct {
	CompileConfig CompileConfig `json:"compile_config"`
	RunConfig     RunConfig     `json:"run_config"`
}

type JudgeResult struct {
	Code     int32  `json:"code"`
	CpuTime  int32  `json:"cpu_time"`
	RealTime int32  `json:"real_time"`
	Memory   int32  `json:"memory"`
	Signal   int32  `json:"signal"`
	ExitCode int32  `json:"exit_code"`
	Error    string `json:"error"`
	Content  string `json:"content"`
	ExePath  string `json:"exe_path"` // 编译result结果才会有
}
