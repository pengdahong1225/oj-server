import request from '@/utils/request'
import type { QueryProblemListParams, SubmitForm, Problem } from '@/types/problem'

// 获取题目列表
export const queryProblemListService = (params: QueryProblemListParams) => {
  return request.get('/problem/list', { params })
}

// 获取题目详情
export const queryProblemDetailService = (id: number) => {
  const params = {
    "problem_id": id
  }
  return request.get(`/problem/detail`, {params})
}

// 获取题目标签列表
export const getProblemTagListService = () => {
  return request.get('/problem/tag_list')
}

// 提交题目
export const submitProblemService = (data: SubmitForm) => {
  return request.post('/problem/submit', data)
}

// 查询提交结果
export const queryResultService = (id: number) => {
  return request.get('/problem/result', { params: { "problem_id": id } })
}

// 更新题目
export const updateProblemService = (data: Problem) => {
  return request.post('/problem/update', data)
}

// 删除题目
export const deleteProblemService = (id: number) => {
  return request.delete(`/problem?problem_id=${id}`)
}