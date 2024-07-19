package judgeService

// JudgeConfig 判题服务的运行参数配置
type JudgeConfig struct {
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

	// 以下字段不属于判题服务的返回字段
	Content string   `json:"content,omitempty"`
	Test    TestCase `json:"test,omitempty"`
}

// 用户状态枚举
const (
	UserStateNormal = iota
	UserStateJudging
)

// “提交”状态枚举，如果没有查询到状态，就意味着最近没有提交题目or题目提交过期了
const (
	UPStateNormal    = iota // 初始状态
	UPStateCompiling        // 正在编译
	UPStateJudging          // 已编译成功，正在判题
	UPStateExited           // 已退出 -> 如何查询到这个状态，就意味着可以查询结果了
)
