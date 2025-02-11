package types

import "github.com/pengdahong1225/oj-server/backend/proto/pb"

type Param struct {
	Uid           int64
	ProblemID     int64
	Code          string // 源代码
	Language      string // 语种
	ProblemConfig *pb.ProblemConfig

	// 编译结果中读取
	FileIds map[string]string // 文件id

	Accepted bool
	Message  string
}

// SandBoxApiForm 沙箱接口表单
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

type Body struct {
	Cmd []SandBoxApiForm `json:"cmd"`
}
