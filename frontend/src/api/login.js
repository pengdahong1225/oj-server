/**
 * 登录相关的请求接口
 */

import request from '@/utils/request'

// 获取图像验证码
export const getPicCode = () => {
  return request.get('/captcha/image')
}

// 获取短信验证码
export const getSmsCode = (mobile, captchaID, captchaValue) => {
  const data = {
    mobile: mobile,
    captchaID: captchaID,
    captchaValue: captchaValue
  }

  return request.post('/captcha/sms', data)
}

// 手机号验证码登录
export const mobileLogin = (mobile, smsCode) => {
  const data = {
    mobile: mobile,
    smsCode: smsCode
  }

  return request.post('/login', data)
}
