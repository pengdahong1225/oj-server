package box

import (
	"bytes"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"io"
	"judge-service/models"
	"os/exec"
	"runtime"
)

// GoSandBox Golang
type GoSandBox struct {
}

func (receiver *GoSandBox) CheckCodeValid(data []byte) error {
	code := string(data)
	for i := 0; i < len(code)-6; i++ {
		if code[i:i+6] == "import" {
			var flag byte
			for i = i + 7; i < len(code); i++ {
				if code[i] == ' ' {
					continue
				}
				flag = code[i]
				break
			}
			if flag == '(' {
				for i = i + 1; i < len(code); i++ {
					if code[i] == ')' {
						break
					}
					if code[i] == '"' {
						t := ""
						for i = i + 1; i < len(code); i++ {
							if code[i] == '"' {
								break
							}
							t += string(code[i])
						}
						if _, ok := models.ValidGolangPackageMap[t]; !ok {
							return errors.Errorf("import %s is not valid", t)
						}
					}
				}
			} else if flag == '"' {
				t := ""
				for i = i + 1; i < len(code); i++ {
					if code[i] == '"' {
						break
					}
					t += string(code[i])
				}
				if _, ok := models.ValidGolangPackageMap[t]; !ok {
					return errors.Errorf("import %s is not valid", t)
				}
			}
		}
	}
	return nil
}

// Run 运行代码
func (receiver *GoSandBox) Run(path string, rsp *models.JudgeBack, testCase models.TestCase) bool {
	cmd := exec.Command("go run", path)
	var out, stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = &out
	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		logrus.Errorf("StdinPipe error: %v", err)
		rsp.Status = models.EN_Status_Internal
		rsp.Tips = err.Error()
		rsp.Output = stderr.String()
		return false
	}
	io.WriteString(stdinPipe, testCase.Input+"\n")

	// 初始内存
	var bm runtime.MemStats
	runtime.ReadMemStats(&bm)

	// 执行
	if err := cmd.Run(); err != nil {
		logrus.Infof("cmd.Run error: %v", err)
		if err.Error() == "exit status 2" {
			rsp.Status = models.EN_Status_CompileError
			rsp.Tips = "Compile Error"
			rsp.Output = stderr.String()
			return false
		}
	}
	// 运行结束内存
	var em runtime.MemStats
	runtime.ReadMemStats(&em)

	// 内存溢出
	if em.Alloc/1024-(bm.Alloc/1024) > uint64(5) {
		rsp.Status = models.EN_Status_MemoryOver
		rsp.Tips = "Memory Over"
		rsp.Output = out.String()
		return false
	}
	// 答案错误
	if testCase.Output != out.String() {
		rsp.Status = models.EN_Status_Faild
		rsp.Tips = "Answer Error"
		rsp.Output = out.String()
		return false
	}
	return true
}
