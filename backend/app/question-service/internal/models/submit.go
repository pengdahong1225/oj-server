package models

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
