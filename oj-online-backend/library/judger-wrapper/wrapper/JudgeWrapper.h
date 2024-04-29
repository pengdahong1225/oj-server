#include <unistd.h>
#include <dlfcn.h>
#include <sys/types.h>
#include <string>
#include <filesystem>
#include <vector>

#include "../common/define.h"

/*
 * JudgeWrapper作为judger-core的外部调用接口类
 * JudgeClient对外不可见
 */
class JudgeWrapper {
public:
    static JudgeResultList
    judge(LangConfig *language_config, std::string &src, int submission_id, const std::string &test_case_json);

private:
    static void writeUtf8ToFile(const std::string &filePath, const std::string &content);
    static std::string readFileContent(const std::filesystem::path &filePath);
    static JudgeResult
    compile(CompileConfig *compile_config, const std::string &src_path, const std::string &work_dir);
    static void initTestCaseEnv(const std::string &work_dir, const std::string &test_case_json);
};
