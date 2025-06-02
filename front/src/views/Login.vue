<template>
  <div class="login-container">
    <div class="login-content">
      <!-- 左侧装饰区域 -->
      <div class="login-left">
        <div class="brand-content">
          <img src="../assets/logo.svg" alt="Logo" class="logo" />
          <h1>Cloudisk 网盘系统</h1>
          <p class="brand-desc">安全、高效、便捷的云存储解决方案</p>
        </div>
      </div>
      
      <!-- 右侧登录表单 -->
      <div class="login-right">
        <div class="login-box">
          <h2>欢迎登录</h2>
          <p class="login-subtitle">请使用您的账号密码登录系统</p>
          <el-form :model="loginForm" :rules="rules" ref="loginFormRef" class="login-form">
            <el-form-item prop="username">
              <el-input 
                v-model="loginForm.username" 
                placeholder="请输入用户名"
                size="large"
              >
                <template #prefix>
                  <el-icon><User /></el-icon>
                </template>
              </el-input>
            </el-form-item>
            <el-form-item prop="password">
              <el-input 
                v-model="loginForm.password" 
                type="password" 
                placeholder="请输入密码"
                size="large"
                show-password
                @keyup.enter="handleLogin"
              >
                <template #prefix>
                  <el-icon><Lock /></el-icon>
                </template>
              </el-input>
            </el-form-item>
            <el-form-item>
              <el-button 
                type="primary" 
                @click="handleLogin" 
                :loading="loading" 
                class="login-button"
                size="large"
              >
                {{ loading ? '登录中...' : '登录' }}
              </el-button>
            </el-form-item>
            <div class="login-options">
              <el-checkbox v-model="rememberMe">记住我</el-checkbox>
              <a href="#" class="forgot-password">忘记密码？</a>
            </div>
            <div class="register-link">
              还没有账号？<router-link to="/register" class="register-btn">立即注册</router-link>
            </div>
          </el-form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { User, Lock } from '@element-plus/icons-vue'
import request from '../utils/request'

const router = useRouter()
const loading = ref(false)
const rememberMe = ref(false)
const loginForm = ref({
  username: '',
  password: ''
})

// 如果之前记住了用户名，自动填充
onMounted(() => {
  const savedUsername = localStorage.getItem('username')
  if (savedUsername) {
    loginForm.value.username = savedUsername
    rememberMe.value = true
  }
})

const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' }
  ]
}

const handleLogin = async () => {
  try {
    loading.value = true
    const formData = new FormData()
    formData.append('username', loginForm.value.username)
    formData.append('password', loginForm.value.password)
    
    const response = await request.post('/user/login', formData)
    
    if (response.code === 200) {
      localStorage.setItem('token', response.data['X-Token'])
      if (rememberMe.value) {
        localStorage.setItem('username', loginForm.value.username)
      } else {
        localStorage.removeItem('username')
      }
      ElMessage.success('登录成功')
      router.push('/files')
    }
  } catch (error) {
    // 错误处理已经在请求拦截器中完成
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  width: 100vw;
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background: #f7fafd;
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
}

.login-content {
  display: flex;
  width: 100%;
  height: 100%;
  background: #fff;
  overflow: hidden;
}

.login-left {
  width: 50%;
  height: 100%;
  background: linear-gradient(135deg, #1e3c72 0%, #2a5298 100%);
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  padding: 40px;
}

.brand-content {
  color: white;
  text-align: center;
}

.logo {
  width: 100px;
  height: 100px;
  margin-bottom: 20px;
}

.brand-content h1 {
  font-size: 32px;
  margin-bottom: 16px;
  font-weight: 600;
}

.brand-desc {
  font-size: 16px;
  opacity: 0.9;
}

.login-right {
  width: 50%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px;
  background: #f7fafd;
}

.login-box {
  width: 100%;
  max-width: 400px;
}

.login-box h2 {
  font-size: 28px;
  color: #1e3c72;
  margin-bottom: 8px;
  font-weight: 600;
}

.login-subtitle {
  color: #666;
  font-size: 14px;
  margin-bottom: 32px;
}

.login-form {
  margin-top: 20px;
}

.login-form :deep(.el-input__wrapper) {
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.05);
  border-radius: 8px;
  padding: 8px 16px;
}

.login-form :deep(.el-input__inner) {
  height: 40px;
}

.login-button {
  width: 100%;
  height: 44px;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 500;
  background: linear-gradient(135deg, #1e3c72 0%, #2a5298 100%);
  border: none;
  transition: transform 0.2s;
}

.login-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(30, 60, 114, 0.2);
}

.login-options {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin: 16px 0;
  font-size: 14px;
}

.forgot-password {
  color: #1e3c72;
  text-decoration: none;
  transition: color 0.2s;
}

.forgot-password:hover {
  color: #2a5298;
}

.register-link {
  text-align: center;
  margin-top: 24px;
  font-size: 14px;
  color: #666;
}

.register-btn {
  color: #1e3c72;
  text-decoration: none;
  font-weight: 500;
  margin-left: 4px;
  transition: color 0.2s;
}

.register-btn:hover {
  color: #2a5298;
}

:deep(.el-checkbox__label) {
  color: #666;
}

:deep(.el-checkbox__input.is-checked .el-checkbox__inner) {
  background-color: #1e3c72;
  border-color: #1e3c72;
}

:deep(.el-checkbox__input.is-checked + .el-checkbox__label) {
  color: #1e3c72;
}

@media screen and (max-width: 1024px) {
  .login-content {
    flex-direction: column;
  }
  
  .login-left, .login-right {
    width: 100%;
    height: 50%;
    padding: 20px;
  }
  
  .login-left {
    padding: 40px 20px;
  }
  
  .login-box {
    padding: 20px;
  }
}
</style> 