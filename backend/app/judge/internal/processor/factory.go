package processor

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"oj-server/app/judge/internal/define"
	"oj-server/module/configManager"
	"strings"
)

func NewProcessor(language string) *BaseProcessor {
	var processor IProcessor

	language = strings.ToLower(language)
	switch language {
	case "c":
		processor = &CProcessor{}
	case "cpp":
		processor = &CProcessor{}
	case "go":
		processor = &GoProcessor{}
	case "python":
		processor = &PyProcessor{}
	default:
		panic("unsupported language")
	}

	// 查询sandbox地址
	var addr string
	for _, item := range configManager.Instance().SandBoxCfg {
		if item.Type == language {
			addr = fmt.Sprintf("http://%s:%d", item.Host, item.Port)
			break
		}
	}
	if addr == "" {
		logrus.Fatalf("target sandbox config not found, language=%s", language)
	}

	return &BaseProcessor{
		impl:       processor,
		sandBoxUrl: addr,
		runResults: make(chan *define.RunResultInChan, 100),
	}
}
