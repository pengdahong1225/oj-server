/**
 * 题目相关的请求接口
 */

import request from '@/utils/request'

// 获取题库列表
export const getProblemList = () => {
  return request.get('problemSet')
}
