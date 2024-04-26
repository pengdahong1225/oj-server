package judger

import "judge-service/models"

type JudgeWrapper struct {
}

func (receiver *JudgeWrapper) Judge(cfg *models.LangConfig, src string, submissionId int32, testCaseJson string) []models.JudgeResult {

	return nil
}

func (receiver *JudgeWrapper) writeUtf8ToFile(filePath *string, content *string) {

}
func (receiver *JudgeWrapper) readFileContent(filePath *string) string {
	return ""
}

func (receiver *JudgeWrapper) compile(cfg *models.CompileConfig, srcPath string, workDir string) []models.JudgeResult {
}

func (receiver *JudgeWrapper) initTestCaseEnv(workDir string, testCaseJson string) {

}
