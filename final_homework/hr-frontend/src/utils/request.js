import axios from 'axios'
import { ElMessage } from 'element-plus'

// 创建 axios 实例
const service = axios.create({
  baseURL: 'http://localhost:8080/api/v1', // 我们的 Gin 网关地址
  timeout: 60000 // 请求超时时间 60 秒（AI 接口需要较长响应时间）
})

// 请求拦截器：自动给每个请求带上 Token
service.interceptors.request.use(
  config => {
    // 从 localStorage 中获取我们登录时存进去的 Token
    const token = localStorage.getItem('token')
    if (token) {
      config.headers['Authorization'] = `Bearer ${token}`
    }
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

// 响应拦截器：统一处理后端的报错
service.interceptors.response.use(
  response => {
    const res = response.data
    // 如果 code 不是 0，说明业务报错了，直接弹出错误提示
    if (res.code !== 0) {
      ElMessage.error(res.msg || 'Error')
      return Promise.reject(new Error(res.msg || 'Error'))
    }
    return res
  },
  error => {
    ElMessage.error(error.message || '网络请求异常')
    return Promise.reject(error)
  }
)

export default service