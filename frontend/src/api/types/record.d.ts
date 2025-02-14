namespace API {
    // 查询记录列表参数
    type QueryRecordListParams = {
        page?: number;
        page_size?: number;
        status?: string;
    }
    // 记录
    type Record = {
        id: number;
        create_at: number;
        uid: number;
        problem_id: number;
        problem_name: string;
        status: string;
        results: Result[];
        code: string;
        lang: string;
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
