namespace API {
    // 查询记录列表参数
    type QueryRecordListParams = {
        page?: number;
        page_size?: number;
    }
    // 记录
    type Record = {
        id: number;
        uid: number;
        problem_id: number;
        created_at: number;
        problem_name?: string;
        user_name?: string;
        // 结果信息
        accepted: boolean;
        message: string;
        status: string;
        lang: string;
        clock: number;
        memory: number;
        code?: string;
        results?: Result[];
    }
    type Result = {
        status: string;
        err_msg: string;
        exitStatus: number;
        time: number;
        memory: number;
        runTime: number;
        content: string;
        testCase: any;
    }
}
