namespace API {
    // 查询题目列表参数
    type QueryProblemListParams = {
        page: number;
        page_size: number;
        keyword: string;
        tag: string;
        uid?: number;
    }
    // 提交题目表单
    type SubmitForm = {
        problem_id: number;
        title: string;
        lang: string;
        code: string;
    }
    
    type Problem = {
        problem_title: string;
        description: string;
        level: number;
        tags: string[];
    
        problem_id: number;
        status?: number;
        create_at?: string;
        create_by?: number;
        config?: string;
    }
    type CreateProblemForm = {
        title: string;
        level: number;
        tags: string[];
        description: string;
    }
    type UpdateProblemForm = {
        id: number;
        title: string;
        level: number;
        tags: string[];
        description: string;
    }
    
    type ProblemConfig = {
        testCases: TestCase[];
        compileLimit: Limit;
        runLimit: Limit;
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
