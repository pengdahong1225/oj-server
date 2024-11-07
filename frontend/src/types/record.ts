export interface QueryRecordListParams{
    page: number;
    page_size: number;
    status: string;
}

export interface Record {
    created_at: string;
    created_by: number;
    id: number;
    lang: string;
    memory: number;
    problem_id: number;
    run_time: number;
    status: number;
    title: string;
    uid: number;
}