import request from '@/utils/request'
import type { QueryProblemListParams, SubmitForm, Problem } from '@/types/problem'

export const queryProblemListService = (params: QueryProblemListParams) => {
  return request.get('/problemSet', { params })
}

export const queryProblemDetailService = (id: number) => {
  const params = {
    "problemID": id
  }
  return request.get(`/problem/detail`, {params})
}

export const submitProblemService = (data: SubmitForm) => {
  return request.post('/problem/submit', data)
}

export const updateProblemService = (data: Problem) => {
  return request.post('/problem/update', data)
}