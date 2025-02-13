import request from '@/utils/request'

/**
 * 获取顶层评论列表接口
 * GET /comment/root_list
 */
export async function getRootCommentListService(params: API.QueryRootCommentListParams) {
    return request('/comment/root_list', {
        method: 'GET',
        params: {
            ...params
        }
    })
}
