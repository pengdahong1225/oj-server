export interface QueryRecordListParams{
    page: number;
    page_size: number;
}

export interface Record {
    create_at: string;
    create_by: number;
    id: number;
    lang: string;
    memory: number;
    problem_id: number;
    run_time: number;
    status: string;
    title: string;
    uid: number;
}