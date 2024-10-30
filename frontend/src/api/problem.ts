import request from '@/utils/request'
import type { QueryProblemListParams } from '@/types/problem'

export const queryProblemListService = (params: QueryProblemListParams) => {
  return request.get('/problemSet', { params })
}

export const queryProblemDetailService = (id: number) => {
  const params = {
    "problemID": id
  }
  return request.get(`/problem/detail`, {params})
}
