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

/**
 * 获取二层评论列表接口
 * GET /comment/child_list
 */
export async function getChildCommentListService(params: API.QueryChildCommentListParams) {
    return request('/comment/child_list', {
        method: 'GET',
        params: {
            ...params
        }
    })
}

/**
 * 提交评论接口
 * POST /comment/add
 */
export async function addCommentService(form: API.AddCommentForm) {
    return request('/comment/add', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        data: form
    })
}

/** 
 * 评论点赞接口
 * POST /comment/like
 */
export async function likeCommentService(form: API.CommentLikeForm) {
    return request('/comment/like', {
        method: 'POST',
        data: form,
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded'
        }
    })
}