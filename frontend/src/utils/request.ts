import axios from 'axios'
import { useUserStore } from '@/stores'
import router from '@/router'

// const baseURL = 'http://192.168.201.128/api'
const baseURL = 'http://localhost:9000'

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
    if (resp.data.code === 0) {
      return resp
    }
    if (resp.data.message) {
      ElMessage.error(resp.data.message)
    } else {
      ElMessage.error('服务异常')
    }
    return Promise.reject(resp.data)
  },
  async (err) => {
    const originalRequest = err.config;
    // 处理 401 错误（Token 过期）
    if (err.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;
      // 如果正在刷新token，就将请求挂起，等待刷新完成
      if (isRefreshing) {
        return new Promise((resolve) => {
          pendingRequests.push(() => resolve(instance(originalRequest)));
        });
      }

      isRefreshing = true;
      try {
        const newToken = await refreshToken();
        isRefreshing = false;
        
        // 更新token
        userStore.setToken(newToken);

        // 更新原始请求的 Authorization 头
        originalRequest.headers.Authorization = `Bearer ${newToken}`;
        
        // 重试原始请求并执行挂起的请求
        const retryResponse = await instance(originalRequest);
        pendingRequests.forEach(cb => cb());
        pendingRequests = [];
        
        return retryResponse;
      } catch (err){
        // 刷新失败，跳转登录页
        userStore.clearToken();
        return Promise.reject(err);
      }
    }

    // 其他错误
    ElMessage.error(err.response?.data?.message || '服务异常')
    return Promise.reject(err)
  }
)

export default instance
export { baseURL }
