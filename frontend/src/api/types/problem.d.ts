namespace API {
    // 查询题目列表参数（请求 query，字段名不变）
    type QueryProblemListParams = {
        page: number;
        page_size: number;
        keyword: string;
        tag: string;
        uid?: number;
    }
    // 提交题目表单（请求体，字段名以后端 model.SubmitRequest 为准，仍是 snake_case）
    type SubmitForm = {
        problem_id: number;
        title: string;
        lang: string;
        code: string;
    }
    // 判题结果摘要（响应 QueryJudgeResultResponse，camelCase）
    type JudgeResultAbstract = {
        accepted: boolean;
        message: string;
    }

    // 题目基础信息（响应 protojson camelCase；int64 字段在读取处用 Number() 转回数字）
    type Problem = {
        id: number;
        createAt?: number;
        updateAt?: number;
        deleteAt?: string;
        title: string;
        description: string;
        level: number;
        tags: string[];
        status: number;
    }
    // 创建题目表单（请求体）
    type CreateProblemForm = {
        problem_title: string;
        level: number;
        tags: string[];
        description: string;
    }
    // 更新题目表单（请求体）
    type UpdateProblemForm = {
        problem_id: number;
        problem_title: string;
        level: number;
        tags: string[];
        description: string;
    }

    // 题目配置（响应 GetProblemConfigResponse.problem_config，camelCase）
    type ProblemConfig = {
        id: number;
        compileLimit: Limit;
        runLimit: Limit;
        testCases: TestCase[];
    }
    type TestCase = {
        input: string;
        output: string;
    }
    type Limit = {
        cpuLimit: number;
        clockLimit: number;
        memoryLimit: number;
        stackLimit: number;
        procLimit: number;
    }
}
