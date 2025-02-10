import request from '@/utils/request'

/**
 * 获取题目列表
 * GET /problem/list
 */
export async function queryProblemListService(params: API.QueryProblemListParams) {
    return request('/problem/list', {
        method: 'GET',
        params: {
            ...params
        }
    })
}
/**
 * 获取题目详情
 * GET /problem/detail
 */
export async function queryProblemDetailService(id: number) {
    return request('/problem/detail', {
        method: 'GET',
        params: { "problem_id": id }
    })
}

/**
 * 获取题目标签列表
 * GET /problem/tag_list
 */
export async function getProblemTagListService(){
    return request('/problem/tag_list', {
        method: 'GET'
    })
}

/**
 * 提交题目
 * POST /problem/submit
 */
export async function submitProblemService(form: API.SubmitForm) {
    return request('/problem/submit', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        data: form
    })
}

/**
 * 查询提交结果
 * GET /problem/result
 */
export async function queryResultService(id: number) {
    return request('/problem/result', {
        method: 'GET',
        params: { "problem_id": id }
    })
}

/**
 * 更新题目
 * POST /problem/update
 */
export async function updateProblemService(data: API.Problem) {
    return request('/problem/update', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        data: data
    })
}

/**
 * 删除题目
 * DELETE /problem
 */
export async function deleteProblemService(id: number) {
    return request('/problem', {
        method: 'DELETE',
        params: { problem_id: id }
    })
}
