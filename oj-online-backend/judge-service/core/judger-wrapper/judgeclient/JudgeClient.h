//
// Created by Messi on 24-4-22.
//

#ifndef CORE_WRAPPER_JUDGECLIENT_H
#define CORE_WRAPPER_JUDGECLIENT_H

#include <string>
#include "../common/define.h"
#include "../common/json/json.h"


/*
 * core调用实例，可用对象池优化，避免频繁创建销毁
 * 初始化之后先读取测试用例
 */
class JudgeClient {
    struct RunConfig _run_config;
    std::string _exe_path;
    std::string _test_case_dir;
    std::string _log_path;
    std::string _work_dir;
    std::string _io_mode;

    Json::Value _test_case_info; // 从info文件加载的测试用例信息
    bool _is_test_case_info_loaded = false;

public:
    explicit JudgeClient(struct RunConfig run_config, std::string exe_path, std::string log_path,
                         std::string test_case_dir,
                         std::string work_dir,
                         std::string io_mode = standardIO);
    void judge(JudgeResultList &resultList);
private:
    void _load_test_case_info();
    JudgeResult _judge_one(int test_case_file_id);
    void _compare_output(int test_case_file_id, const std::string &user_output_file);

private:
    static std::string readFileContent(const std::filesystem::path &filePath);
};


#endif //CORE_WRAPPER_JUDGECLIENT_H
