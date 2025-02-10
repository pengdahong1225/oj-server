namespace API {
    // 查询记录列表参数
    type QueryRecordListParams = {
        page?: number;
        page_size?: number;
        status?: string;
    }
    // 记录
    type Record = {
        created_at?: string;
        created_by?: number;
        id?: number;
        lang?: string;
        memory?: number;
        problem_id?: number;
        run_time?: number;
        status?: number;
        title?: string;
        uid?: number;
    }
}
