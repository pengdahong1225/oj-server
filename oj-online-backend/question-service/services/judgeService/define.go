package judgeService

// 部分编译和运行配置
type CompileConfig struct {
	CpuLimit    int64 `json:"cpu_limit,omitempty"`    // CPU时间限制，单位纳秒
	ClockLimit  int64 `json:"clock_limit,omitempty"`  // 等待时间限制，单位纳秒 （通常为 cpuLimit 两倍）
	MemoryLimit int64 `json:"memory_limit,omitempty"` // 内存限制，单位 byte
	ProcLimit   int64 `json:"proc_limit,omitempty"`   // 线程数量限制
}
type RunConfig struct {
	CpuLimit    int64 `json:"cpu_limit,omitempty"`    // CPU时间限制，单位纳秒
	ClockLimit  int64 `json:"clock_limit,omitempty"`  // 等待时间限制，单位纳秒 （通常为 cpuLimit 两倍）
	MemoryLimit int64 `json:"memory_limit,omitempty"` // 内存限制，单位 byte
	ProcLimit   int64 `json:"proc_limit,omitempty"`   // 线程数量限制
}

type TestCase struct {
	Input  string `json:"input,omitempty"`
	Output string `json:"output,omitempty"`
}

type Result struct {
	Status     string            `json:"status,omitempty"`
	Error      string            `json:"error,omitempty"`      // 详细错误信息
	ExitStatus int64             `json:"exitStatus,omitempty"` // 程序返回值
	Time       int64             `json:"time,omitempty"`       // 程序运行 CPU 时间，单位纳秒
	Memory     int64             `json:"memory,omitempty"`     // 程序运行内存，单位 byte
	RunTime    int64             `json:"runTime,omitempty"`    // 程序运行现实时间，单位纳秒
	Files      map[string]string `json:"files,omitempty"`      // copyOut 和 pipeCollector 指定的文件内容
	FileIds    map[string]string `json:"fileIds,omitempty"`    // copyFileCached 指定的文件 id
	FileError  []string          `json:"fileError,omitempty"`  // 文件错误详细信息
}

type ResultInCache struct {
	Content string    `json:"content,omitempty"`
	Results []*Result `json:"results,omitempty"`
}

// 用户状态枚举
const (
	UserStateNormal = iota
	UserStateJudging
)

// 用户提交题目状态枚举
const (
	UPStateNormal    = iota
	UPStateCompiling // 正在编译
	UPStateRunning   // 正在运行
	UPStateJudging   // 正在判题
	UPStateExited    // 已退出
)
