package processor

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"oj-server/app/judge/internal/biz"
	"oj-server/app/judge/internal/define"
	"oj-server/module/configManager"
	"strings"
)

func NewProcessor(language string, uc *biz.JudgeUseCase) *BaseProcessor {
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
		logrus.Errorf("language not supported, language=%s", language)
		return nil
	}

	// 查询sandbox地址
	var addr string
	for _, item := range configManager.AppConf.SandBoxCfg {
		if item.Type == language {
			addr = fmt.Sprintf("http://%s:%d", item.Host, item.Port)
			break
		}
	}
	if addr == "" {
		logrus.Fatalf("target sandbox config not found, language=%s", language)
	}

	return &BaseProcessor{
		uc:         uc,
		impl:       processor,
		sandBoxUrl: addr,
		runResults: make(chan *define.RunResultInChan, 100),
	}
}
