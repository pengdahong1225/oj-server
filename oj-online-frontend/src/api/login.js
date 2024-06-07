/**
 * 登录相关的请求接口
 */

import request from '@/utils/request'

// 获取图像验证码
export const getPicCode = () => {
  return request.get('/captcha/image')
}

// 获取短信验证码
export const getSmsCode = (mobile, picKey, picCode) => {
  return request.post('/captcha/sendSmsCaptcha', {
    form: {
      mobile: mobile,
      captchaKey: picKey,
      captchaCode: picCode
    }
  })
}

// 手机号验证码登录
export const phoneLogin = (mobile, smsCode) => {
  return request.post('/passport/login', {
    form: {
      isParty: false,
      partyData: {},
      mobile: mobile,
      smsCode: smsCode
    }
  })
}
