<template>
  <div class="register-container">
    <div class="register-content">
      <!-- 左侧品牌区 -->
      <div class="register-left">
        <div class="brand-content">
          <img src="../assets/logo.svg" alt="Logo" class="logo" />
          <h1>Cloudisk 网盘系统</h1>
          <p class="brand-desc">安全、高效、便捷的云存储解决方案</p>
        </div>
      </div>
      <!-- 右侧注册表单区 -->
      <div class="register-right">
        <div class="register-box">
          <h2>欢迎注册</h2>
          <p class="register-subtitle">请填写信息注册新账号</p>
          <el-form :model="registerForm" :rules="rules" ref="registerFormRef" class="register-form">
            <el-form-item prop="username">
              <el-input 
                v-model="registerForm.username" 
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
                v-model="registerForm.password" 
                type="password" 
                placeholder="请输入密码"
                size="large"
                show-password
              >
                <template #prefix>
                  <el-icon><Lock /></el-icon>
                </template>
              </el-input>
            </el-form-item>
            <el-form-item prop="confirmPassword">
              <el-input 
                v-model="registerForm.confirmPassword" 
                type="password" 
                placeholder="请再次输入密码"
                size="large"
                show-password
              >
                <template #prefix>
                  <el-icon><Lock /></el-icon>
                </template>
              </el-input>
            </el-form-item>
            <el-form-item>
              <el-button 
                type="primary" 
                @click="handleRegister" 
                :loading="loading" 
                class="register-button"
                size="large"
              >
                {{ loading ? '注册中...' : '注册' }}
              </el-button>
            </el-form-item>
            <div class="login-link">
              已有账号？<router-link to="/login" class="login-btn">立即登录</router-link>
            </div>
          </el-form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { User, Lock } from '@element-plus/icons-vue'
import request from '../utils/request'

const router = useRouter()
const loading = ref(false)
const registerForm = ref({
  username: '',
  password: '',
  confirmPassword: ''
})

const validatePass = (rule, value, callback) => {
  if (!value) {
    callback(new Error('请输入密码'))
  } else if (registerForm.value.confirmPassword && value !== registerForm.value.confirmPassword) {
    callback(new Error('两次输入密码不一致'))
  } else {
    callback()
  }
}

const validatePass2 = (rule, value, callback) => {
  if (!value) {
    callback(new Error('请再次输入密码'))
  } else if (value !== registerForm.value.password) {
    callback(new Error('两次输入密码不一致'))
  } else {
    callback()
  }
}

const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '长度在 3 到 20 个字符', trigger: 'blur' }
  ],
  password: [
    { required: true, validator: validatePass, trigger: 'blur' },
    { min: 6, message: '密码长度不能小于6位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, validator: validatePass2, trigger: 'blur' }
  ]
}

const handleRegister = async () => {
  try {
    loading.value = true
    const formData = new FormData()
    formData.append('username', registerForm.value.username)
    formData.append('password', registerForm.value.password)
    
    const response = await request.post('/user/register', formData)
    
    if (response.code === 200) {
      ElMessage.success('注册成功')
      router.push('/login')
    }
  } catch (error) {
    // 错误处理已经在请求拦截器中完成
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.register-container {
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

.register-content {
  display: flex;
  width: 100%;
  height: 100%;
  background: #fff;
  overflow: hidden;
}

.register-left {
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

.register-right {
  width: 50%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px;
  background: #f7fafd;
}

.register-box {
  width: 100%;
  max-width: 400px;
}

.register-box h2 {
  font-size: 28px;
  color: #1e3c72;
  margin-bottom: 8px;
  font-weight: 600;
}

.register-subtitle {
  color: #666;
  font-size: 14px;
  margin-bottom: 32px;
}

.register-form {
  margin-top: 20px;
}

.register-form :deep(.el-input__wrapper) {
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.05);
  border-radius: 8px;
  padding: 8px 16px;
}

.register-form :deep(.el-input__inner) {
  height: 40px;
}

.register-button {
  width: 100%;
  height: 44px;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 500;
  background: linear-gradient(135deg, #1e3c72 0%, #2a5298 100%);
  border: none;
  transition: transform 0.2s;
}

.register-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(30, 60, 114, 0.2);
}

.login-link {
  text-align: center;
  margin-top: 24px;
  font-size: 14px;
  color: #666;
}

.login-btn {
  color: #1e3c72;
  text-decoration: none;
  font-weight: 500;
  margin-left: 4px;
  transition: color 0.2s;
}

.login-btn:hover {
  color: #2a5298;
}

@media screen and (max-width: 1024px) {
  .register-content {
    flex-direction: column;
  }
  
  .register-left, .register-right {
    width: 100%;
    height: 50%;
    padding: 20px;
  }
  
  .register-left {
    padding: 40px 20px;
  }
  
  .register-box {
    padding: 20px;
  }
}
</style> 