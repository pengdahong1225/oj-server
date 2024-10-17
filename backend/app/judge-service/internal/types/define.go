package types

type Param struct {
	Uid       int64
	ProblemID int64
	Code      string // 源代码
	Language  string // 语种

	// 从redis读取
	CompileLimit Limit
	RunLimit     Limit
	TestCases    []TestCase

	// 编译结果中读取
	FileIds map[string]string // 文件id
}

// ProblemHotData redis 题目热点数据模型
type ProblemHotData struct {
	CompileLimit string `json:"compile_limit"`
	RunLimit     string `json:"run_limit"`
	TestCases    string `json:"test_cases"`
}

type StaticConfig struct {
	Argc          []string                     // 程序命令行参数
	Env           []string                     // 环境变量
	Files         []map[string]any             // 指定 标准输入、标准输出、标准错误的文件
	CopyIn        map[string]map[string]string // 在执行程序之前复制进容器的文件列表
	CopyOut       []string                     // 在执行程序后从容器文件系统中复制出来的文件列表（内容）
	CopyOutCached []string                     // 和 copyOut 相同，不过文件不返回内容，而是返回一个对应文件 ID
}

type Limit struct {
	CpuLimit    int64 `json:"cpuLimit,omitempty"`    // CPU时间限制，单位纳秒
	ClockLimit  int64 `json:"clockLimit,omitempty"`  // 等待时间限制，单位纳秒 （通常为 cpuLimit 两倍）
	MemoryLimit int64 `json:"memoryLimit,omitempty"` // 内存限制，单位 byte
	StackLimit  int64 `json:"stackLimit,omitempty"`  // 栈内存限制，单位 byte
	ProcLimit   int64 `json:"procLimit,omitempty"`   // 线程数量限制
}
type TestCase struct {
	Input  string `json:"input,omitempty"`
	Output string `json:"output,omitempty"`
}

type SubmitResult struct {
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
