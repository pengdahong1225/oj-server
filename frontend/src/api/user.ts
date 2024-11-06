import request from '@/utils/request'
import type { RegisterForm, LoginForm } from '@/types/user'

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
export const querySubmitRecordService = (id: number) => {
    const params = {
        "problem_id": id
    }
    return request.get('/user/submit_record', { params })
}
