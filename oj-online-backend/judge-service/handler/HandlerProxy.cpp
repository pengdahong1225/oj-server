//
// Created by Messi on 24-4-29.
//

#include "HandlerProxy.h"
#include "common/define.h"
#include "wrapper/JudgeWrapper.h"

SSJudgeResponse HandlerProxy::handle(SSJudgeRequest &request) {
    SSJudgeResponse response;
    response.set_session_id(request.session_id());
    // 一系列校验
    if (request.code().empty() || request.language().empty() || request.test_case_json().empty()) {
        auto &result = *response.add_result_list();
        result.set_result(-1);
        return response;
    }
    // 业务逻辑
    {
        CompileConfig compile_config;
        {
            compile_config.src_name = "main.cpp";
            compile_config.exe_name = "main";
            compile_config.max_cpu_time = 1000;
            compile_config.max_real_time = 2000;
            compile_config.max_memory = 128 * 1024 * 1024;
            compile_config.compiler_exe = "/usr/bin/g++";
//            compile_config.compile_args = "-std=c++17 -O2 -w -fmax-errors=3 -DONLINE_JUDGE -lm -s -static";
        }
        RunConfig run_config;
        {
            run_config.seccomp_rule_name = "c_cpp";
            run_config.memory_limit_check_only = 1;
            run_config.max_cpu_time = 1000;
            run_config.max_real_time = 2000;
            run_config.max_memory = 128 * 1024 * 1024;
        }
        LangConfig lang_config;
        lang_config.compile_config = compile_config;
        lang_config.run_config = run_config;

        // judge
        auto result_list = JudgeWrapper::judge(&lang_config, const_cast<std::string &>(request.code()),
                                               request.submit_id(),
                                               request.test_case_json());

        for (const auto &item: result_list) {
            auto &result = *response.add_result_list();
            result.set_result(item.code);
            result.set_cpu_time(item.cpu_time);
            result.set_real_time(item.real_time);
            result.set_memory(item.memory);
            result.set_signal(item.signal);
            result.set_exit_code(item.exit_code);
            result.set_error(item.error);
            result.set_content(item.content);
        }
    }

    return response;
}
