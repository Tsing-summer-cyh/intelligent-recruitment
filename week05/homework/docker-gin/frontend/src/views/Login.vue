<template>
  <div class="login-container">
    <div class="login-box">
      <h2>{{ isLogin ? '登录智能单词本' : '注册新账号' }}</h2>
      <div class="form-group">
        <input v-model="form.username" type="text" placeholder="请输入用户名" />
      </div>
      <div class="form-group">
        <input v-model="form.password" type="password" placeholder="请输入密码" />
      </div>
      <div class="form-actions">
        <button @click="handleSubmit" class="primary-btn">
          {{ isLogin ? '登 录' : '注 册' }}
        </button>
      </div>
      <div class="toggle-text" @click="isLogin = !isLogin">
        {{ isLogin ? '没有账号？点击注册' : '已有账号？点击登录' }}
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import request from '../api/request';

const router = useRouter();
const isLogin = ref(true);
const form = ref({ username: '', password: '' });

const handleSubmit = async () => {
  if (!form.value.username || !form.value.password) {
    alert("用户名和密码不能为空！");
    return;
  }
  
  try {
    if (isLogin.value) {
      // 登录请求
      const res = await request.post('/login', form.value);
      localStorage.setItem('jwt_token', res.token); // 保存 Token
      alert("登录成功");
      router.push('/dashboard'); // 跳转到主页
    } else {
      // 注册请求
      await request.post('/register', form.value);
      alert("注册成功，请登录");
      isLogin.value = true; // 切换回登录模式
    }
  } catch (error) {
    alert(error.response?.data?.error || "操作失败，请重试");
  }
};
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  /* 酷炫的蓝紫渐变背景 */
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); 
}
.login-box {
  background: rgba(255, 255, 255, 0.95);
  padding: 40px;
  border-radius: 16px;
  box-shadow: 0 10px 25px rgba(0,0,0,0.2);
  width: 320px;
  text-align: center;
  backdrop-filter: blur(10px); /* 毛玻璃效果 */
}
h2 {
  margin-top: 0; color: #1f2937; font-weight: 600; margin-bottom: 24px;
}
.form-group input {
  width: 100%; padding: 12px 16px; margin-bottom: 16px;
  border: 1px solid #d1d5db; border-radius: 8px;
  box-sizing: border-box; font-size: 15px;
  transition: all 0.3s ease;
}
.form-group input:focus {
  outline: none; border-color: #4f46e5;
  box-shadow: 0 0 0 3px rgba(79, 70, 229, 0.2); /* 聚焦时的光晕 */
}
.primary-btn {
  width: 100%; padding: 12px;
  background-color: #4f46e5; color: white;
  border: none; border-radius: 8px; font-size: 16px; font-weight: 600;
  cursor: pointer; transition: all 0.2s; margin-top: 10px;
}
.primary-btn:hover { background-color: #4338ca; transform: translateY(-1px); box-shadow: 0 4px 12px rgba(79, 70, 229, 0.3); }
.primary-btn:active { transform: translateY(0); }
.toggle-text {
  margin-top: 20px; color: #4f46e5; cursor: pointer; font-size: 14px; transition: color 0.2s;
}
.toggle-text:hover { color: #3730a3; text-decoration: underline; }
</style>