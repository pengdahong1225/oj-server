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
    // 新增公告表单（请求体）
    type AddNoticeForm = {
        title: string;
        content: string;
    }
    // 更新公告表单（请求体）
    type UpdateNoticeForm = {
        id: number;
        title: string;
        content: string;
    }
}
