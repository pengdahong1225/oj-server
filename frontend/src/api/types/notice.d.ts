namespace API {
    type QueryNoticeListParams = {
        page: number;
        page_size: number;
    }
    type Notice = {
        id: number;
        title: string;
        createAt: string;
        content: string;
        status: number;
        createBy: number;
        createByName?: string;
    }
}
