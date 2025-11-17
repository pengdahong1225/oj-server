import axios from 'axios'
import { useUserStore } from '@/stores'
import router from '@/router'

const baseURL = 'http://localhost:8080/api/v1'

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
    console.log(resp)
    
    if (resp.data.code === 0) {
        return resp.data.data.access_token;
    }
    return null;
  } catch (err) {
    // 刷新失败，跳转到登录页
    ElMessage.error('登录已过期，请重新登录');
    router.push('/login');
    throw err;
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
    // 处理 401 Token 过期
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
        const newToken = await refreshToken()
        isRefreshing = false

        // 设置 token
        userStore.setToken(newToken)

        // 更新 header
        originalRequest.headers.Authorization = `Bearer ${newToken}`

        // 先执行挂起请求
        pendingRequests.forEach((cb) => cb())
        pendingRequests = []

        // 重试当前请求
        return instance(originalRequest)
      } catch (refreshErr) {
        isRefreshing = false
        pendingRequests = []

        // 刷新失败 → 清空 token → 跳登录
        userStore.clearToken()
        return Promise.reject(refreshErr)
      }
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
