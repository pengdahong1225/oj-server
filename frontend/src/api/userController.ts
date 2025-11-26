import request from '@/utils/request'

/**
 * 注册接口
 * POST /user/register
 */
export async function userRegisterService(form: API.RegisterForm) {
    return request('/user/register', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        data: form
    })
}

/**
 * 登录接口
 * POST /user/login
 */
export async function userLoginService(form: API.LoginForm) {
    return request('/user/login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        data: form
    })
}
/**
 * 获取用户信息
 * GET /user/profile
 */
export async function getUserInfoService(uid: number) {
    return request('/user/profile', {
        method: 'GET',
        params: { uid },
    })
}

/**
 * 查询用户提交记录列表
 * GET /user/record_list
 */
export async function queryRecordListService(params: API.QueryRecordListParams) {
    return request('/record/record_list', {
        method: 'GET',
        params: {
            ...params
        }
    })
}
/**
 * 查询用户提交记录
 * GET /user/record
 */
export async function queryRecordService(id: number) {
    return request('/record/record', {
        method: 'GET',
        params: { id }
    })
}


/**
 * 获取用户列表
 * GET /user/list
 */
export async function queryUserListService(params: API.QueryUserListParams) {
    return request('/user/list', {
        method: 'GET',
        params: {
            ...params
        }
    })
}