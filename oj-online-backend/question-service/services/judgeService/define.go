package judgeService

type CompileConfig struct {
	CpuLimit    int64  `json:"cpu_limit,omitempty"`    // CPU时间限制，单位纳秒
	ClockLimit  int64  `json:"clock_limit,omitempty"`  // 等待时间限制，单位纳秒 （通常为 cpuLimit 两倍）
	MemoryLimit int64  `json:"memory_limit,omitempty"` // 内存限制，单位 byte
	ProcLimit   int64  `json:"proc_limit,omitempty"`   // 线程数量限制
	Content     string `json:"content,omitempty"`      // 源代码
}

type RunConfig struct {
	CpuLimit    int64  `json:"cpu_limit,omitempty"`    // CPU时间限制，单位纳秒
	ClockLimit  int64  `json:"clock_limit,omitempty"`  // 等待时间限制，单位纳秒 （通常为 cpuLimit 两倍）
	MemoryLimit int64  `json:"memory_limit,omitempty"` // 内存限制，单位 byte
	ProcLimit   int64  `json:"proc_limit,omitempty"`   // 线程数量限制
	StdIn       string `json:"std_in,omitempty"`       // 输入
	FileId      string `json:"file_id,omitempty"`      // 文件id（这个缓存文件的 ID 来自上一个请求返回的 fileIds）
}
