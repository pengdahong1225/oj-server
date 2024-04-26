#include "JudgeWrapper.h"
#include <fstream>
#include <iostream>
#include "../common/json/json.h"
#include "../judgeclient/JudgeClient.h"
extern "C" {
#include "../core/runner.h"
}

/*
 * language_config：配置
 * src：提交的源代码
 * submission_id：唯一提交id
 * test_case_json：测试用例，需要解析
 */
JudgeResultList JudgeWrapper::judge(LangConfig *language_config, std::string &src, int submission_id,
                                    const std::string &test_case_json) {
    JudgeResultList resultList;
    auto &compile_config = language_config->compile_config;
    auto &run_config = language_config->run_config;
    std::string work_dir = judger_dir + std::to_string(submission_id);
    std::string log_path = work_dir + "/" + "judge.log";
    std::string test_case_dir = work_dir + "/" + "test_case";

    // TODO write source code into file
    // 源代码文件路径：/app/judger/{submission_id}/{src_name}
    std::string src_path = work_dir + "/" + compile_config.src_name;
    JudgeResult init_result{};
    try {
        writeUtf8ToFile(src_path, src);
    } catch (const std::exception &e) {
        sprintf(init_result.content.data(), "Failed to write source code into file[%s],err=%s",
                src_path.c_str(),
                e.what());
        resultList.emplace_back(init_result);
        return resultList;
    }

    // TODO compile source code, return exe file path
    JudgeResult compile_result = compile(&compile_config, src_path, work_dir);
    if (compile_result.code != SUCCESS || compile_result.exe_path.empty()) {
        resultList.emplace_back(compile_result);
        return resultList;
    }

    // TODO 初始化测试环境
    JudgeResult parse_result{};
    try {
        initTestCaseEnv(work_dir, test_case_json);
    } catch (const std::exception &e) {
        sprintf(parse_result.content.data(), "Failed to parse test case json,err=%s", e.what());
        resultList.emplace_back(parse_result);
        return resultList;
    }

    // TODO judge
    auto judger_client = new JudgeClient(run_config, compile_result.exe_path, log_path, test_case_dir, work_dir);
    judger_client->judge(resultList);
    delete judger_client;
    return resultList;
}

void JudgeWrapper::writeUtf8ToFile(const std::string &filePath, const std::string &content) {
    std::ofstream file(filePath, std::ios_base::out | std::ios_base::trunc);
    if (!file.is_open()) {
        throw std::runtime_error("Failed to open file: " + filePath);
    }

    file << content;
    file.flush();
    file.close();
}

JudgeResult
JudgeWrapper::compile(CompileConfig *compile_config, const std::string &src_path, const std::string &work_dir) {
    JudgeResult ret{};

    std::string exe_path = work_dir + "/" + compile_config->exe_name;
    std::string compiler_out = work_dir + "/" + "compiler.out";
    std::string log_path = work_dir + "/" + "compile.log";

    // 构造编译命令
    // /usr/bin/g++ -O2 -w -fmax-errors=3 -std=c++11 {src_path} -lm -o {exe_path}
    std::string cmd =
            compile_config->compiler_exe + " " + compile_config->compile_args + " " + src_path + " -lm -o " + exe_path;

    // 编译
    struct config cfg{
            // 限制
            .max_cpu_time = static_cast<int>(compile_config->max_cpu_time),
            .max_real_time = compile_config->max_real_time,
            .max_memory = compile_config->max_memory,
            .max_stack = 128 * 1024 * 1024,
            .max_process_number = UNLIMITED,
            .max_output_size = 20 * 1024 * 1024,
            // 执行参数
            .exe_path = compile_config->compiler_exe.data(), // 编译器
            .input_path = const_cast<char *>(src_path.data()),
            .output_path = compiler_out.data(),
            .error_path = compiler_out.data(),
            .args = {compile_config->compile_args.data()}, // 编译参数
            .env = {nullptr},
            .log_path = log_path.data(),
            .seccomp_rule_name = nullptr,
            .uid = 0,
            .gid = 0
    };
    struct result result{};
    run(&cfg, &result);
    ret.code = result.result;
    ret.cpu_time = result.cpu_time;
    ret.real_time = result.real_time;
    ret.memory = result.memory;
    ret.signal = result.signal;
    ret.exit_code = result.exit_code;
    ret.error = result.error;

    if (result.result != SUCCESS) {
        ret.content = "Compile failed";
        return ret;
    }
    else {
        if (!std::filesystem::exists(compiler_out)) {
            sprintf(ret.content.data(), "Failed to read compiler_out file,err = %s",
                    "compiler_out file not exists");
            return ret;
        }
        std::string content;
        try {
            content = readFileContent(compiler_out);
        } catch (const std::exception &e) {
            sprintf(ret.content.data(), "Failed to read compiler_out file,err = %s", e.what());
            return ret;
        }
        ret.content = content;
        ret.exe_path = exe_path;
        return ret;
    }
}

std::string JudgeWrapper::readFileContent(const std::filesystem::path &filePath) {
    std::ifstream file(filePath);
    if (!file.is_open()) {
        throw std::runtime_error("Failed to open file: " + filePath.string());
    }

    std::string content((std::istreambuf_iterator<char>(file)), std::istreambuf_iterator<char>());
    file.close();

    return content;
}

/*
 * 初始化测试用例环境
 * 1.创建test_case目录
 * 2.解析test_case_json
 * 3.创建文件
 */
void JudgeWrapper::initTestCaseEnv(const std::string &work_dir, const std::string &test_case_json) {
    // work_dir = /app/judger/{submission_id}
    // 创建test_case目录
    if (std::filesystem::create_directory(work_dir + "/test_case")) {
        throw std::runtime_error("Failed to create test_case directory: " + work_dir + "/test_case");
    }

    // 解析test_case_json
    Json::Value root; // 顶层
    Json::Reader reader;
    if (!reader.parse(test_case_json, root) || root.empty()) {
        throw std::runtime_error("Failed to parse test_case_json: " + test_case_json);
    }

    // info
    std::string info_path = work_dir + "/test_case/info.json";
    try {
        writeUtf8ToFile(info_path, root.get("info", "").asString());
    } catch (const std::exception &e) {
        throw e; // 抛到上层
    }

    // input
    if (!root.isMember("input")) {
        throw std::runtime_error("Failed to parse test_case_json: " + test_case_json);
    }
    if (!root["input"].isArray()) {
        throw std::runtime_error("Failed to parse test_case_json: " + test_case_json);
    }
    Json::Value input_array = root["input"];
    for (const auto &item: input_array) {
        std::string input_path = work_dir + "/test_case/" + item.get("name", "").asString();
        std::string content = item.get("content", "").asString();
        try {
            writeUtf8ToFile(input_path, content);
        } catch (const std::exception &e) {
            throw e; // 抛到上层
        }
    }
    // output
    if (!root.isMember("output")) {
        throw std::runtime_error("Failed to parse test_case_json: " + test_case_json);
    }
    if (!root["output"].isArray()) {
        throw std::runtime_error("Failed to parse test_case_json: " + test_case_json);
    }
    Json::Value output_array = root["output"];
    for (const auto &item: output_array) {
        std::string input_path = work_dir + "/test_case/" + item.get("name", "").asString();
        std::string content = item.get("content", "").asString();
        try {
            writeUtf8ToFile(input_path, content);
        } catch (const std::exception &e) {
            throw e; // 抛到上层
        }
    }
}
