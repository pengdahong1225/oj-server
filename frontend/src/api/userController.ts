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
 * 创建用户接口（管理员）
 * POST /user/add
 */
export async function createUserService(form: API.CreateUserForm) {
    return request('/user/add', {
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
 * 退出登录
 * POST /user/logout
 */
export async function userLogoutService() {
    return request('/user/logout', {
        method: 'POST',
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
 * GET /record/list
 */
export async function queryRecordListService(params: API.QueryRecordListParams) {
    return request('/record/list', {
        method: 'GET',
        params: {
            ...params
        }
    })
}
/**
 * 查询用户提交记录详情
 * GET /record/detail
 */
export async function queryRecordService(id: number) {
    return request('/record/detail', {
        method: 'GET',
        params: { record_id: id }
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