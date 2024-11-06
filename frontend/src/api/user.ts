import request from '@/utils/request'
import type { RegisterForm, LoginForm } from '@/types/user'
import type { QueryRecordListParams } from '@/types/record'

// 注册接口
export const userRegisterService = (form: RegisterForm) => {
    return request.post('/user/register', form)
}
// 登录接口
export const userLoginService = (form: LoginForm) => {
    return request.post('/user/login', form)
}
// 获取用户信息
export const getUserInfoService = (uid: number) => {
    return request.get(`/user/profile?uid=${uid}`)
}

// 查询用户提交记录
export const queryRecordListService = (params: QueryRecordListParams) => {
    return request.get('/user/record_list', { params })
}
export const queryRecordService = (id: number) => {
    return request.get(`/user/record?id=${id}`)
}
