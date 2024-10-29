import request from '@/utils/request'
import type { RegisterForm, LoginForm, UPSSParams } from '@/types/user'

// 注册接口
export const userRegisterService = (form: RegisterForm) => {
    return request.post('/register', form)
}
// 登录接口
export const userLoginService = (form: LoginForm) => {
    return request.post('/login', form)
}
// 查询题目集用户解决状态
export const upssService = (params: UPSSParams) => {
    return request.get('/user/upss', { params })
}

// 获取用户信息
export const getUserInfoService = (uid: number) => {
    return request.get(`/user/profile?uid=${uid}`)
}