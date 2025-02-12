import request from '@/utils/request'

/**
 * 获取评论列表接口
 * GET /comment/query
 */
export async function getCommentListService(params: API.QueryCommentListParams) {
    return request('/comment/query', {
        method: 'GET',
        params: {
            ...params
        }
    })
}
