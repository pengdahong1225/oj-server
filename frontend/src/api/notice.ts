import request from '@/utils/request'
import type { QueryNoticeListParams } from '@/types/notice'

export const queryNoticeListService = (params: QueryNoticeListParams) => {
    return request.get('/notice/noticeList', { params })
}
