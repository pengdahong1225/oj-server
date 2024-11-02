import request from '@/utils/request'
import type { RegisterForm, LoginForm } from '@/types/user'

// 注册接口
export const userRegisterService = (form: RegisterForm) => {
    return request.post('/register', form)
}
// 登录接口
export const userLoginService = (form: LoginForm) => {
    return request.post('/login', form)
}
// 获取用户信息
export const getUserInfoService = (uid: number) => {
    return request.get(`/user/profile?uid=${uid}`)
}