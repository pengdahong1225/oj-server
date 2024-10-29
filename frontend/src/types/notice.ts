export interface QueryNoticeListParams {
    page: number;
    page_size: number;
    keyword: string;
}
export interface Notice {
    id: number;
    title: string;
    create_at: string;
    content: string;
    status: number;
    create_by: number;
}