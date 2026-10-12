namespace API {
    type QueryNoticeListParams = {
        page: number;
        page_size: number;
        keyword: string;
    }
    // 公告（响应 Notice，camelCase）
    type Notice = {
        id: number;
        title: string;
        createAt: string;
        content: string;
        status: number;
        createBy: number;
    }
}
