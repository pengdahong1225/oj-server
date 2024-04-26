package box

import (
	"github.com/pkg/errors"
	"judge-service/models"
)

// SandBox 沙箱接口
type SandBox interface {
	CheckCodeValid([]byte) error
	Run(string, *models.JudgeBack, models.TestCase) bool
}

func NewSandBox(lang string) (SandBox, error) {
	switch lang {
	case "cpp":
		return new(CppSandBox), nil
	case "go":
		return new(GoSandBox), nil
	default:
		return nil, errors.Errorf("unsupported language: %s", lang)
	}
}
