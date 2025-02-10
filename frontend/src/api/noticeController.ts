import request from '@/utils/request'

// 查询公告列表
export async function queryNoticeListService(params: API.QueryNoticeListParams) {
    return request('/notice/list', {
        method: 'GET',
        params: {
            ...params
        }
    })
}
