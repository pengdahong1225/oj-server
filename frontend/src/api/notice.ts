import request from '@/utils/request'
import type { QueryNoticeListParams } from '@/types/notice'

// 查询公告列表
export const queryNoticeListService = (params: QueryNoticeListParams) => {
    return request.get('/notice/list', { params })
}
