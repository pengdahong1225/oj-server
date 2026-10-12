namespace API {
    // 查询记录列表参数（请求 query）
    type QueryRecordListParams = {
        page?: number;
        page_size?: number;
    }
    // 提交记录（响应 SubmitRecord，camelCase）
    type Record = {
        id: number;
        uid: number;
        problemId: number;
        problemName: string;
        message: string;
        results?: Result[];
        code?: string;
        lang: string;
        createdAt: number;
        problemLevel: number;
        userName: string;
        accepted: boolean;
        // 以下字段新版 proto SubmitRecord 未携带，前端暂留空（待后端补充或从 results 汇总）
        clock?: number;
        memory?: number;
    }
    // 判题结果子项（响应 JudgeResultItem，camelCase）
    type Result = {
        status: string;
        errMsg: string;
        exitStatus: number;
        time: number;
        memory: number;
        runTime: number;
        content: string;
        testCase: any;
    }
}
