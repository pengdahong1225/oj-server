/**
 * 用户相关的请求接口
 */

import request from '@/utils/request'

// 获取用户信息
export const getUserProfile = (uid) => {
  console.log(uid)
  const params = {
    uid: uid
  }
  return request.get('/user/profile', { params: params })
}

// 获取用户解题信息
export const getUserSolvedList = () => {
  return request.get('/user/solvedList')
}
