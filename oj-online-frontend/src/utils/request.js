import axios from 'axios'
import store from '@/store'
import { Loading } from 'element-ui' // 引入element-ui的loading组件

/**
 * 创建axios实例，将来对创建出来的axios实例，进行自定义配置
 * 好处：不会污染原始的axios实例
 * 自定义配置：请求/响应 拦截器
 */

// 创建axios实例
const instance = axios.create({
  baseURL: 'http://smart-shop.itheima.net/index.php?s=/api',
  timeout: 5000
})

let loading
function startLoading () { // 使用Element loading-start 方法
  loading = Loading.service({
    lock: true,
    text: '加载中……'
  })
}
function endLoading () { // 使用Element loading-close 方法
  loading.close()
}

// 添加请求拦截器
instance.interceptors.request.use(function (config) {
  // 在发送请求之前做些什么
  // 开启loading 禁止背景点击 以服务的方式调用的全屏 Loading 是单例的
  startLoading()
  // 配置headers，有token就添加token
  config.headers.platform = 'H5'
  const token = store.getters.token
  if (token) {
    config.headers['Access-Token'] = token
  }

  return config
}, function (error) {
  // 对请求错误做些什么
  return Promise.reject(error)
})

// 添加响应拦截器
instance.interceptors.response.use(function (response) {
  // 对响应数据做点什么(默认axios的响应数据结构会多包装一层data)
  const res = response.data
  console.log(res)
  if (res.status !== 200) {
    endLoading()
    return Promise.reject(res.message)
  } else {
    endLoading()
  }

  return res
}, function (error) {
  // 超出 2xx 范围的状态码都会触发该函数。
  // 对响应错误做点什么
  return Promise.reject(error)
})

/**
   * 导出配置好的axios实例
   */
export default instance
