import request from '@/utils/request'
import type { QueryProblemListParams } from '@/types/problem'

export const queryProblemListService = (params: QueryProblemListParams) => {
  return request.get('/problemSet', { params })
}

export const queryProblemDetailService = (id: number) => {
  return request.get(`/problem/${id}`)
}