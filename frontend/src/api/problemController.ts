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
export async function queryResultService(task_id: number) {
    return request('/problem/result', {
        method: 'GET',
        params: { "task_id": task_id }
    })
}

/**
 * 创建题目
 * POST /problem/add
 */
export async function addProblemService(data: API.CreateProblemForm) {
    return request('/problem/add', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        data: data
    })
}

/**
 * 更新题目
 * POST /problem/update
 */
export async function updateProblemService(data: API.UpdateProblemForm) {
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

/**
 * 发布题目
 * POST /problem/publish
 */
export async function publishProblemService(id: number) {
    return request('/problem/publish', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded'
        },
        data: {
            problem_id: id
        }
    })
}

/**
 * 隐藏题目
 * POST /problem/hide
 */
export async function hideProblemService(id: number) {
    return request('/problem/hide', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded'
        },
        data: {
            problem_id: id
        }
    })
}

/**
 * 上传题目配置
 * POST /problem/upload_config
 */
export async function uploadProblemConfigService(data: FormData) {
    return request('/problem/upload_config', {
        method: 'POST',
        data: data
    })
}