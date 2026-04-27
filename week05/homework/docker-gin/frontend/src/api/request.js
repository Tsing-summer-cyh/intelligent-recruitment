// frontend/src/api/request.js
import axios from 'axios';

const request = axios.create({
  baseURL: '/api', // 配合 vite.config.ts 的 proxy 解决开发环境跨域
  timeout: 10000,
});

// 请求拦截器：自动携带 Token
request.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('jwt_token');
    if (token) {
      config.headers['Authorization'] = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// 响应拦截器：统一处理错误
request.interceptors.response.use(
  (response) => response.data,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('jwt_token');
      window.location.href = '/login'; // Token 失效直接踢回登录页
    }
    return Promise.reject(error);
  }
);

export default request;