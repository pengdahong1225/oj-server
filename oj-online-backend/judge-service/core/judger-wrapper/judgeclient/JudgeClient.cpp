//
// Created by Messi on 24-4-22.
//

#include <fstream>
#include <iostream>
#include <filesystem>

#include "JudgeClient.h"

extern "C" {
#include "../core/runner.h"
}

JudgeClient::JudgeClient(struct RunConfig run_config, std::string exe_path, std::string log_path,
                         std::string test_case_dir,
                         std::string work_dir,
                         std::string io_mode) : _run_config(std::move(run_config)),
                                                _exe_path(std::move(exe_path)),
                                                _log_path(std::move(log_path)),
                                                _test_case_dir(std::move(test_case_dir)),
                                                _work_dir(std::move(work_dir)),
                                                _io_mode(std::move(io_mode)) {
    // 加载测试用例info
    _load_test_case_info();
}

void JudgeClient::judge(JudgeResultList &resultList) {
    // 循环运行测试用例
    auto &test_cases_arr = _test_case_info["test_cases"];
    for (int i = 0; i < test_cases_arr.size(); i++) {
        auto result = _judge_one(i + 1);
        resultList.emplace_back(result);
    }
}

void JudgeClient::_load_test_case_info() {
    // 打开info文件并读取到_test_case_info中
    std::ifstream info_file(_test_case_dir + "/info.json");
    std::stringstream file_contents_stream;
    std::string line;

    if (!info_file) {
        std::cerr << "Error: Unable to open info file " << std::endl;
        return;
    }
    file_contents_stream << info_file.rdbuf(); // 使用rdbuf()直接读取文件内容
    info_file.close();
    std::string file_contents = file_contents_stream.str();

    // json解析
    Json::CharReaderBuilder builder;
    Json::CharReader *reader = builder.newCharReader();
    std::string errors;
    bool parsingSuccessful = reader->parse(file_contents.c_str(), file_contents.c_str() + file_contents.size(),
                                           &_test_case_info,
                                           &errors);
    delete reader;
    if (!parsingSuccessful) {
        std::cerr << "Failed to parse JSON: " << errors << std::endl;
        return;
    }

    _is_test_case_info_loaded = true;
}

// 运行一个测试用例
JudgeResult JudgeClient::_judge_one(int test_case_file_id) {
    JudgeResult ret{};

    if (!_is_test_case_info_loaded) {
        ret.content = "Error: Test case info not loaded";
        return ret;
    }

    // 读取id对应的case
    auto &test_cases = _test_case_info["test_cases"];
    std::string in_file = _test_case_dir + "/" + test_cases[test_case_file_id]["input_name"].asString();

    // 创建运行需要的文件
    std::string real_user_output_file;
    std::string user_output_file;
    if (_io_mode == standardIO) {
        real_user_output_file = user_output_file = _work_dir + "/" + std::to_string(test_case_file_id) + ".out";
    }
    else {
        ret.content = "Error: _io_mode";
        return ret;
    }

    // 构建命令
    std::string cmd = _exe_path;

    // 执行
    struct config cfg{
            // 限制
            .max_cpu_time = _run_config.max_cpu_time,
            .max_real_time = _run_config.max_real_time,
            .max_memory = _run_config.max_memory,
            .max_stack = 128 * 1024 * 1024,
            .max_process_number = UNLIMITED,
            .max_output_size = std::max(1024 * 1024 * 16,
                                        test_cases[test_case_file_id].get("output_size", 0).asInt() * 2),
            // 执行参数
            .exe_path =cmd.data(),
            .input_path = in_file.data(),
            .output_path = real_user_output_file.data(),
            .error_path = real_user_output_file.data(),
            .args = {nullptr},
            .env = {nullptr},
            .log_path = _log_path.data(),
            .seccomp_rule_name = _run_config.seccomp_rule_name.data(),
            .uid = 0,
            .gid = 0
    };
    struct result result{};
    run(&cfg, &result);

    // if progress exited normally, then we should check output result
    ret.code = result.result;
    ret.cpu_time = result.cpu_time;
    ret.real_time = result.real_time;
    ret.memory = result.memory;
    ret.signal = result.signal;
    ret.exit_code = result.exit_code;
    ret.error = result.error;
    if (result.result != SUCCESS) {
        ret.content = "Error: progress exited normally";
        return ret;
    }
    else {
        if (!std::filesystem::exists(user_output_file)) {
            ret.content = "Error: user output file not found";
            return ret;
        }
        try {
            _compare_output(test_case_file_id, user_output_file);
        } catch (const std::exception &e) {
            sprintf(ret.content.data(), "output failed Error: %s", e.what());
            return ret;
        }
        ret.content = "Success";
        return ret;
    }
}

void JudgeClient::_compare_output(int test_case_file_id, const std::string &user_output_file) {
    // 读取用户执行输出文件
    std::string user_content;
    try {
        user_content = readFileContent(user_output_file);
    } catch (const std::exception &e) {
        throw e; // 抛到上层
    }

    // 读取测试用例的输出文件
    std::string case_out_content;
    auto &test_cases = _test_case_info["test_cases"];
    std::string case_out_file = _test_case_dir + "/" + test_cases[test_case_file_id].get("output_name", "").asString();
    try {
        case_out_content = readFileContent(case_out_file);
    } catch (const std::exception &e) {
        throw e; // 抛到上层
    }

    // 比较
    if (user_content.empty()) {
        throw std::runtime_error("Error: user output file is empty");
    }
    if (case_out_content.empty()) {
        throw std::runtime_error("Error: case output file is empty");
    }
    if (user_content != case_out_content) {
        throw std::runtime_error("Error: user output file is not equal to case output file");
    }
}

std::string JudgeClient::readFileContent(const std::filesystem::path &filePath) {
    std::ifstream file(filePath);
    if (!file.is_open()) {
        throw std::runtime_error("Failed to open file: " + filePath.string());
    }

    std::string content((std::istreambuf_iterator<char>(file)), std::istreambuf_iterator<char>());
    file.close();

    return content;
}
