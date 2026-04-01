import axios from 'axios';
import { message } from 'antd';

const request = axios.create({
  baseURL: '', // 依靠 Vite 代理
  timeout: 10000,
});

// 请求拦截器：自动携带 Token
request.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// 响应拦截器：统一处理错误和解包数据
request.interceptors.response.use(
  (response) => {
    const res = response.data;
    // 如果返回的不是标准的 JSON 格式（如静态文件），直接返回
    if (response.config.responseType === 'blob') return response.data;
    
    if (res.code !== 0) {
      message.error(res.msg || '请求失败');
      return Promise.reject(new Error(res.msg));
    }
    return res.data; // 直接返回 data 部分
  },
  (error) => {
    if (error.response?.status === 401) {
      message.error('登录已过期，请重新登录');
      localStorage.removeItem('token');
      window.location.href = '/login';
    } else {
      message.error(error.response?.data?.msg || '网络错误');
    }
    return Promise.reject(error);
  }
);

export default request;