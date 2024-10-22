import request from '@/utils/request'
import type { RegisterParams } from '@/types/user'

// 注册接口
export const userRegisterService = ({ username, password, repassword }: RegisterParams) => {
    return request.post('/api/reg', { username, password, repassword })
}
// 登录接口
export const userLoginService = (username: string, password: string) => {
    return request.post('/api/login', { username, password })
}
// 获取用户信息
export const userGetInfoService = () => {
    return request.get('/my/userinfo')
}
