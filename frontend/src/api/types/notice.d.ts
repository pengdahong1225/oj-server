namespace API {
    type QueryNoticeListParams = {
        page: number;
        page_size: number;
        keyword: string;
    }
    type Notice = {
        id: number;
        title: string;
        create_at: string;
        content: string;
        status: number;
        create_by: number;
    }
}