/**
 * 题目相关的请求接口
 */

import request from '@/utils/request'

// 获取题库列表
export const getProblemList = () => {
  return request.get('problemSet')
}

// 获取题目详情
export const getProblemDetail = (id) => {
  const params = {
    problemID: id
  }
  return request.get('problem/detail', { params: params }).catch(err => {
    console.log(err)
  })
}

// 提交代码
export const submitCode = (obj) => {
  const data = {
    problem_id: obj.problem_id,
    title: obj.title,
    lang: obj.lang,
    code: obj.code
  }
  return request.post('problem/submit', data)
}
