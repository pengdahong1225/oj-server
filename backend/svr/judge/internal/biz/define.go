package biz

import "oj-server/pkg/proto/pb"

// 判题任务上下文参数
type Param struct {
	Code          string // 源代码
	Language      string // 语种
	ProblemConfig *pb.ProblemConfig

	FileIds  map[string]string // 文件id, 从编译结果中读取
	Accepted bool              // 是否通过
	Message  string
	TaskId   string
}

// 沙箱接口表单
type SandBoxApiForm struct {
	// 资源限制
	CpuLimit          int64  `json:"cpuLimit"`          // CPU时间限制，单位纳秒
	ClockLimit        int64  `json:"clockLimit"`        // 等待时间限制，单位纳秒 （通常为 cpuLimit 两倍）
	MemoryLimit       int64  `json:"memoryLimit"`       // 内存限制，单位 byte
	StackLimit        int64  `json:"stackLimit"`        // 栈空间限制，单位 byte
	ProcLimit         int64  `json:"procLimit"`         // 线程数量限制
	CpuRateLimit      int64  `json:"cpuRateLimit"`      // 仅 Linux，CPU 使用率限制，1000 等于单核 100%
	CpuSetLimit       string `json:"cpuSetLimit"`       // 仅 Linux，限制 CPU 使用，使用方式和 cpuset cgroup 相同 （例如，`0` 表示限制仅使用第一个核）
	StrictMemoryLimit bool   `json:"strictMemoryLimit"` // deprecated: 使用 dataSegmentLimit （这个选项依然有效）
	DataSegmentLimit  bool   `json:"dataSegmentLimit"`  // 仅linux，开启 rlimit 堆空间限制（如果不使用cgroup默认开启）
	AddressSpaceLimit bool   `json:"addressSpaceLimit"` // 仅linux，开启 rlimit 虚拟内存空间限制（非常严格，在所以申请时触发限制）

	// 程序命令行参数
	Args []string `json:"args"`
	// 环境变量
	Env []string `json:"env"`

	// 指定 标准输入、标准输出和标准错误的文件
	Files []map[string]any `json:"files"`
	// 开启 TTY （需要保证标准输出和标准错误为同一文件）同时需要指定 TERM 环境变量 （例如 TERM=xterm）
	Tty bool `json:"tty"`

	// 在执行程序之前复制进容器的文件列表
	CopyIn map[string]map[string]string `json:"copyIn"`
	// 在执行程序后从容器文件系统中复制出来的文件列表
	CopyOut []string `json:"copyOut"`
	// 和 copyOut 相同，不过文件不返回内容，而是返回一个对应文件 ID ，内容可以通过 /file/:fileId 接口下载
	CopyOutCached []string `json:"copyOutCached"`
	// 指定 copyOut 复制文件大小限制，单位 byte
	CopyOutMax int64 `json:"copyOutMax"`
}

type SandBoxApiBody struct {
	Cmd []SandBoxApiForm `json:"cmd"`
}

const (
	Accepted            = "Accepted"              // 正常情况
	MemoryLimitExceeded = "Memory Limit Exceeded" // 内存超限
	TimeLimitExceeded   = "Time Limit Exceeded"   // 时间超限
	OutputLimitExceeded = "Output Limit Exceeded" // 输出超限
	FileErr             = "File Error"            // 文件错误
	NonzeroExitStatus   = "Nonzero Exit Status"   // 非 0 退出值
	Signalled           = "Signalled"             // 进程被信号终止
	InternalError       = "Internal Error"        // 内部错误

	WrongAnswer = "Wrong Answer" //结果错误
)

type SandBoxApiResponse struct {
	Status     string `json:"status,omitempty"`
	ErrMsg     string `json:"error,omitempty"`      // 详细错误信息
	ExitStatus int64  `json:"exitStatus,omitempty"` // 程序返回值
	Time       int64  `json:"time,omitempty"`       // 程序运行 CPU 时间，单位纳秒
	Memory     int64  `json:"memory,omitempty"`     // 程序运行内存，单位 byte
	RunTime    int64  `json:"runTime,omitempty"`    // 程序运行现实时间，单位纳秒

	// copyOut 和 pipeCollector 指定的文件内容
	Files map[string]string `json:"files,omitempty"`
	// copyFileCached 指定的文件 id
	FileIds map[string]string `json:"fileIds,omitempty"`
	// 文件错误详细信息
	FileErrs []FileError `json:"fileError,omitempty"`
}
type FileError struct {
	Name    string `json:"name,omitempty"`    // 错误文件名称
	Type    string `json:"type,omitempty"`    // 错误代码
	Message string `json:"message,omitempty"` // 错误信息
}

type RunResultInChan struct {
	Result *SandBoxApiResponse
	Case   *pb.TestCase
}
