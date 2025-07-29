import axios from 'axios'
import { useUserStore } from '@/stores'
import router from '@/router'

// const baseURL = 'http://192.168.201.128/api'
const baseURL = 'http://localhost:9000'

const useStore = useUserStore()
let isRefreshing = false;
let requests: Array<(token: string | null) => void> = [];

function processQueue(error: any, token: string | null = null) {
  requests.forEach(callback => {
    if (error) {
      callback(error);
    } else {
      callback(token);
    }
  });
  requests = [];
}
function refreshToken() {
  return new Promise<string>((resolve, reject) => {
    instance.post('/user/refresh_token').then(res => {
      if (res.data.code === 0) {
        resolve(res.data.token);
      } else {
        reject(res.data.message);
      }
    }).catch(error => {
      reject(error);
    });
  });
}

/**
 * axios
 * 请求拦截器
 * 响应拦截器
 */
const instance = axios.create({
  // 基地址，超时时间
  baseURL,
  timeout: 10000
})
instance.interceptors.request.use(
  (config) => {
    // 携带token
    if (useStore.token) {
      config.headers.Authorization = useStore.token
      config.headers['token'] = useStore.token
    }
    return config
  },
  (err) => Promise.reject(err)
)
instance.interceptors.response.use(
  (res) => {
    // 摘取核心响应数据
    if (res.data.code === 0) {
      return res
    }
    // 处理业务失败, 给错误提示，抛出错误
    ElMessage.error(res.data.message || '服务异常')
    return Promise.reject(res.data)
  },
  (err) => {
    // 错误的特殊情况 => 401 权限不足 或 token 过期
    // if (err.response?.status === 401) {
    //   router.push('/login')
    // }

    const { config, response } = err;
    const originalRequest = config;
    if (response && response.status === 401) {
      if (!isRefreshing) {
        isRefreshing = true;
        return refreshToken().then(newToken => {
          useStore.token = newToken;
          originalRequest.headers['token'] = newToken;
          processQueue(null, newToken);
          return instance(originalRequest);
        }).catch(err => {
          processQueue(err, null);
          useStore.clearUserInfo()
          router.push('/login');
          return Promise.reject(err);
        }).finally(() => {
          isRefreshing = false;
        });
      } else {
        return new Promise((resolve, reject) => {
          requests.push((token) => {
            originalRequest.headers['token'] = token;
            resolve(instance(originalRequest));
          });
        });
      }
    }
    
    ElMessage.error(err.response.data.message || '服务异常')
    return Promise.reject(err)
  }
)

export default instance
export { baseURL }
