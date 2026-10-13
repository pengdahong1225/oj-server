import axios from 'axios'
import { useUserStore } from '@/stores'
import { showLoginDialog } from '@/utils/loginDialog'

const baseURL = 'http://localhost:9000/api/v1'

/**
 * axios
 */
const instance = axios.create({
  // 基地址，超时时间
  baseURL,
  timeout: 10000,
  withCredentials: true,
})

const userStore = useUserStore()
// 是否正在刷新 Token 的标记（防止并发请求重复刷新）
let isRefreshing = false;
// 被挂起的请求队列
let pendingRequests: Array<() => void> = [];

/**
 * 刷新 Token 的函数
 */
const refreshToken = async () => {
  try {
    const resp = await instance.post('/user/refresh_token', null, {
      withCredentials: true, // 确保携带 Cookie
    });
    
    if (resp.data.code === 0) {
        return resp.data.data.access_token;
    }
    return null;
  } catch (err) {
    // 刷新失败，静默处理（由外层统一弹「请先登录」）
    return null;
  }
};

/**
 * 请求拦截器
 */
instance.interceptors.request.use(
  (config) => {
    // 携带token
    if (userStore.token) {
      config.headers.Authorization = `Bearer ${userStore.token}`;
    }
    return config
  },
  (err) => {
    Promise.reject(err)
  }
)

/**
 * 响应拦截器
 */
instance.interceptors.response.use(
  (resp) => {
    const data = resp.data

    // 业务成功
    if (data?.code === 0) {
      return resp
    }

    // HTTP 成功，但业务失败
    if (data?.message) {
      ElMessage.error(data.message)
    } else {
      ElMessage.error('服务异常')
    }
    return Promise.reject(data)
  },
  // HTTP失败
  async (err) => {
    const originalRequest = err.config

    // ---------------------------
    // 处理 401：未登录 / token 过期
    // ---------------------------
    if (err.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true

      // 若正在刷新 token，则挂起请求
      if (isRefreshing) {
        return new Promise((resolve) => {
          pendingRequests.push(() => resolve(instance(originalRequest)))
        })
      }

      // 开始刷新
      isRefreshing = true
      try {
        // 删除过期token
        userStore.clearToken()

        // 刷新
        const newToken = await refreshToken()
        isRefreshing = false

        if (newToken) {
          // 刷新成功 → 设置新token → 重试
          userStore.setToken(newToken)
          originalRequest.headers.Authorization = `Bearer ${newToken}`
          pendingRequests.forEach((cb) => cb())
          pendingRequests = []
          return instance(originalRequest)
        }

        // 刷新失败 → 需要登录
        pendingRequests = []
        userStore.clearToken()
        ElMessage.warning('请先登录!')
        showLoginDialog()
        return Promise.reject(err)
      } catch (refreshErr) {
        isRefreshing = false
        pendingRequests = []
        userStore.clearToken()
        ElMessage.warning('请先登录!')
        showLoginDialog()
        return Promise.reject(refreshErr)
      }
    }

    // refresh_token 请求自身的错误，静默处理（不弹提示）
    if (originalRequest?.url?.includes('/user/refresh_token')) {
      return Promise.reject(err)
    }

    // ---------------------------
    // 非 401 错误
    // ---------------------------

    const msg =
      err.response?.data?.message ||
      err.response?.data?.msg ||
      err.message ||
      '服务异常'

    ElMessage.error(msg)

    return Promise.reject(err)
  }
)

export default instance
export { baseURL }
