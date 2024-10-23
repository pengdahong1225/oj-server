import request from '@/utils/request'
import type { QueryProblemListParams } from '@/types/problem'

export const queryProblemListService = (params: QueryProblemListParams) => {
  return request.get('/problemSet', { params })
}