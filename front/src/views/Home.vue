<template>
  <div class="home-container">
    <el-container>
      <el-header>
        <div class="header-content">
          <h2>Cloudisk 网盘系统</h2>
          <el-button type="danger" @click="handleLogout">退出登录</el-button>
        </div>
      </el-header>
      <el-main>
        <h3>欢迎使用 Cloudisk 网盘系统</h3>
        <!-- 这里后续会添加文件列表等内容 -->
      </el-main>
    </el-container>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'

const router = useRouter()

const handleLogout = () => {
  localStorage.removeItem('token')
  ElMessage.success('退出成功')
  router.push('/login')
}
</script>

<style scoped>
.home-container {
  height: 100vh;
}

.el-header {
  background-color: #409EFF;
  color: white;
  line-height: 60px;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.el-main {
  padding: 20px;
  background-color: #f5f7fa;
}
</style> 