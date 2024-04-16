package logic

import (
	"github.com/pkg/errors"
	impl2 "judge-service/internal/logic/impl"
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
		return new(impl2.CppSandBox), nil
	case "go":
		return new(impl2.GoSandBox), nil
	default:
		return nil, errors.Errorf("unsupported language: %s", lang)
	}
}
