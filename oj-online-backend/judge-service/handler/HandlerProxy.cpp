//
// Created by Messi on 24-4-29.
//

#include "HandlerProxy.h"
#include "judger-wrapper/common/define.h"
#include "judger-wrapper/wrapper/JudgeWrapper.h"

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

        }
        RunConfig run_config;
        {
            run_config.seccomp_rule_name = "c_cpp";
            run_config.memory_limit_check_only = 1;

        }
        LangConfig lang_config;
        lang_config.compile_config = compile_config;
        lang_config.run_config = run_config;
        // 构造submission_id

        // judge
        auto result_list = JudgeWrapper::judge(&lang_config, const_cast<std::string &>(request.code()), 0,
                                               request.test_case_json());

        for (const auto &item: result_list){
            
        }
    }

    return response;
}
