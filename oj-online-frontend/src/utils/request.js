import axios from 'axios'
import store from '@/store'
import { Loading, Message } from 'element-ui' // 引入element-ui的组件

/**
 * 创建axios实例，将来对创建出来的axios实例，进行自定义配置
 * 好处：不会污染原始的axios实例
 * 自定义配置：请求/响应 拦截器
 */

/**
 * 创建axios实例
 * 配置基地址和超时时间
 */
const instance = axios.create({
  // baseURL: 'http://smart-shop.itheima.net/index.php?s=/api',
  baseURL: 'http://127.0.0.1:9010',
  timeout: 5000
})

/**
 * element loading
 */
function startLoading () {
  Loading.service({
    // lock: true,
    text: 'Loading',
    spinner: 'el-icon-loading',
    background: 'rgba(0, 0, 0, 0.8)'
  })
}
function endLoading () {
  Loading.service().close()
}

// 添加请求拦截器
instance.interceptors.request.use(function (config) {
  // 在发送请求之前做些什么
  // 开启loading 禁止背景点击 以服务的方式调用的全屏 Loading 是单例的
  startLoading()
  // 配置headers，有token就添加token
  // config.headers.platform = 'H5'
  const token = store.getters.token
  if (token) {
    config.headers.token = token
  }

  return config
}, function (error) {
  // 对请求错误做些什么
  Message.error(error.message)
  return Promise.reject(error)
})

// 添加响应拦截器
instance.interceptors.response.use(function (response) {
  // 对响应数据做点什么(默认axios的响应数据结构会多包装一层data)
  endLoading()
  if (response.data.code !== 200) {
    Message.error(response.data.message)
    return Promise.reject(response.data.message)
  }
  return response.data
}, function (error) {
  // 超出 2xx 范围的状态码都会触发该函数。
  // 对响应错误做点什么
  endLoading()
  console.log(error)
  if (error.response) {
    Message.error(error.response.data.message)
  } else {
    Message.error(error.message)
  }
  return Promise.reject(error)
})

/**
   * 导出配置好的axios实例
   */
export default instance
