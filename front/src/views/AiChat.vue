<template>
  <div class="file-list-container">
    <div class="navbar">
      <div class="navbar-left">
        <img src="../assets/logo.svg" alt="Logo" class="navbar-logo" />
        <span class="navbar-title">Cloudisk 网盘系统</span>
      </div>
      <div class="navbar-right">
        <span class="navbar-user">{{ username }}</span>
        <el-button type="text" class="navbar-logout" @click="handleLogout">退出</el-button>
      </div>
    </div>
    <div class="main-layout">
      <div class="sidebar">
        <div class="sidebar-menu">
          <div class="sidebar-menu-item" :class="{active: currentView==='files'}" @click="currentView='files'">
            <el-icon><Folder /></el-icon>
            <span>我的文件</span>
          </div>
          <div class="sidebar-menu-item" :class="{active: currentView==='transfer'}" @click="currentView='transfer'">
            <el-icon><UploadFilled /></el-icon>
            <span>传输管理</span>
          </div>
          <div class="sidebar-menu-item" :class="{active: currentView==='download'}" @click="currentView='download'">
            <el-icon><Download /></el-icon>
            <span>下载管理</span>
          </div>
          <div class="sidebar-menu-item" :class="{active: currentView==='recycle'}" @click="currentView='recycle'">
            <el-icon><Delete /></el-icon>
            <span>回收站</span>
          </div>
          <div class="sidebar-menu-item" :class="{active: currentView==='ai-chat'}" @click="currentView='ai-chat'">
            <el-icon><ChatDotRound /></el-icon>
            <span>智能文件助手</span>
          </div>
          <div class="sidebar-menu-item" :class="{active: currentView==='setting'}" @click="currentView='setting'">
            <el-icon><Setting /></el-icon>
            <span>设置</span>
          </div>
        </div>
      </div>
      <div class="main-content">
        <div class="chat-content">
          <!-- 聊天窗口 -->
          <div class="chat-window">
            <!-- 消息列表 -->
            <div class="message-list" ref="messageList">
              <div v-for="(msg, index) in messages" :key="index" 
                   :class="['message', msg.type === 'user' ? 'user-message' : 'ai-message']">
                <div class="message-avatar">
                  <el-avatar :size="40" :icon="msg.type === 'user' ? UserFilled : Service" />
                </div>
                <div class="message-content">
                  <div class="message-text" v-if="msg.type === 'user'">{{ msg.content }}</div>
                  <div class="message-text markdown-body" v-else v-html="renderMarkdown(msg.content)"></div>
                  <div class="message-time">{{ msg.time }}</div>
                </div>
              </div>
            </div>

            <!-- 输入区域 -->
            <div class="input-area">
              <el-input
                v-model="inputMessage"
                type="textarea"
                :rows="3"
                placeholder="请输入您的问题..."
                @keyup.enter.ctrl="sendMessage"
              />
              <div class="input-actions">
                <span class="tip">按 Ctrl + Enter 发送</span>
                <el-button 
                  type="primary" 
                  :loading="loading"
                  @click="sendMessage"
                >
                  发送
                </el-button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { UserFilled, Service, Folder, ChatDotRound, User, UploadFilled, Download, Delete, Setting } from '@element-plus/icons-vue'
import request from '../utils/request'
import { marked } from 'marked'
import DOMPurify from 'dompurify'

const router = useRouter()
const messages = ref([])
const inputMessage = ref('')
const loading = ref(false)
const messageList = ref(null)
const username = ref(localStorage.getItem('username') || '用户')
const currentView = ref('ai-chat')

// 导航方法
const goToFileList = () => {
  router.push('/file-list')
}

const goToTransfer = () => {
  router.push('/transfer')
}

const goToDownload = () => {
  router.push('/download')
}

const goToAiChat = () => {
  router.push('/ai-chat')
}

const goToUserCenter = () => {
  router.push('/user-center')
}

// 渲染 Markdown
const renderMarkdown = (content) => {
  return DOMPurify.sanitize(marked(content))
}

// 发送消息
const sendMessage = async () => {
  if (!inputMessage.value.trim()) {
    ElMessage.warning('请输入消息内容')
    return
  }

  // 添加用户消息
  const userMessage = {
    type: 'user',
    content: inputMessage.value,
    time: new Date().toLocaleTimeString()
  }
  messages.value.push(userMessage)
  
  // 清空输入框
  const message = inputMessage.value
  inputMessage.value = ''
  
  // 滚动到底部
  await nextTick()
  scrollToBottom()
  
  try {
    loading.value = true
    const response = await request.post('/ai/chat', {
      message: message
    })
    
    if (response.code === 200) {
      // 添加 AI 回复
      messages.value.push({
        type: 'ai',
        content: response.data.message,
        time: new Date().toLocaleTimeString()
      })
    } else {
      ElMessage.error(response.msg || '请求失败')
    }
  } catch (error) {
    console.error('发送消息失败:', error)
    ElMessage.error('发送消息失败')
  } finally {
    loading.value = false
    await nextTick()
    scrollToBottom()
  }
}

// 滚动到底部
const scrollToBottom = () => {
  if (messageList.value) {
    messageList.value.scrollTop = messageList.value.scrollHeight
  }
}

// 组件挂载时滚动到底部
onMounted(() => {
  scrollToBottom()
})
</script>

<style scoped>
.file-list-container {
  position: fixed;
  left: 0;
  top: 0;
  width: 100vw;
  height: 100vh;
  background: linear-gradient(135deg, #f5f7fa 0%, #e4e7eb 100%);
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  overflow: hidden;
}

.navbar {
  width: 100vw;
  height: 64px;
  background: #fff;
  box-shadow: 0 2px 8px 0 rgba(30,60,114,0.06);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 48px;
  z-index: 10;
}

.navbar-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.navbar-logo {
  height: 32px;
}

.navbar-title {
  font-size: 18px;
  font-weight: 600;
  color: #1e3c72;
}

.navbar-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.navbar-user {
  font-size: 14px;
  color: #1e3c72;
}

.main-layout {
  display: flex;
  width: 100vw;
  height: calc(100vh - 64px);
}

.sidebar {
  width: 200px;
  height: 100%;
  background: #fff;
  box-shadow: 2px 0 8px 0 rgba(30,60,114,0.04);
  display: flex;
  flex-direction: column;
  align-items: center;
  padding-top: 32px;
}

.sidebar-menu {
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.sidebar-menu-item {
  width: 100%;
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 32px;
  font-size: 16px;
  color: #1e3c72;
  border-radius: 8px 0 0 8px;
  cursor: pointer;
  transition: background 0.2s, color 0.2s;
}

.sidebar-menu-item.active,
.sidebar-menu-item:hover {
  background: linear-gradient(90deg, #e4e7eb 0%, #f5f7fa 100%);
  color: #2a5298;
}

.main-content {
  flex: 1;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: flex-start;
  overflow-y: auto;
}

.chat-content {
  width: 100%;
  max-width: 900px;
  padding: 32px;
}

.chat-window {
  width: 100%;
  max-width: 900px;
  height: calc(100vh - 64px);
  background: #fff;
  border-radius: 16px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
  display: flex;
  flex-direction: column;
}

.message-list {
  flex: 1;
  overflow-y: auto;
  padding: 24px;
}

.message {
  display: flex;
  margin-bottom: 24px;
  animation: fadeIn 0.3s ease;
}

.message-avatar {
  margin-right: 12px;
}

.message-content {
  max-width: 80%;
}

.message-text {
  padding: 12px 16px;
  border-radius: 12px;
  font-size: 14px;
  line-height: 1.6;
  word-break: break-word;
}

.message-time {
  font-size: 12px;
  color: #999;
  margin-top: 4px;
}

.user-message {
  flex-direction: row-reverse;
}

.user-message .message-avatar {
  margin-right: 0;
  margin-left: 12px;
}

.user-message .message-text {
  background: linear-gradient(135deg, #1e3c72 0%, #2a5298 100%);
  color: #fff;
}

.ai-message .message-text {
  background: #f5f7fa;
  color: #1e3c72;
}

/* Markdown 样式 */
.markdown-body {
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif;
  line-height: 1.6;
}

.markdown-body h1,
.markdown-body h2,
.markdown-body h3,
.markdown-body h4,
.markdown-body h5,
.markdown-body h6 {
  margin-top: 24px;
  margin-bottom: 16px;
  font-weight: 600;
  line-height: 1.25;
}

.markdown-body h1 { font-size: 2em; }
.markdown-body h2 { font-size: 1.5em; }
.markdown-body h3 { font-size: 1.25em; }
.markdown-body h4 { font-size: 1em; }
.markdown-body h5 { font-size: 0.875em; }
.markdown-body h6 { font-size: 0.85em; }

.markdown-body p {
  margin-top: 0;
  margin-bottom: 16px;
}

.markdown-body ul,
.markdown-body ol {
  margin-top: 0;
  margin-bottom: 16px;
  padding-left: 2em;
}

.markdown-body li {
  margin-top: 0.25em;
}

.markdown-body code {
  padding: 0.2em 0.4em;
  margin: 0;
  font-size: 85%;
  background-color: rgba(27,31,35,0.05);
  border-radius: 3px;
}

.markdown-body pre {
  margin-top: 0;
  margin-bottom: 16px;
  padding: 16px;
  overflow: auto;
  font-size: 85%;
  line-height: 1.45;
  background-color: #f6f8fa;
  border-radius: 3px;
}

.markdown-body blockquote {
  padding: 0 1em;
  color: #6a737d;
  border-left: 0.25em solid #dfe2e5;
  margin: 0 0 16px 0;
}

.markdown-body strong {
  font-weight: 600;
}

.markdown-body em {
  font-style: italic;
}

.input-area {
  padding: 24px;
  border-top: 1px solid #f0f0f0;
}

.input-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 12px;
}

.tip {
  font-size: 12px;
  color: #999;
}

:deep(.el-textarea__inner) {
  border-radius: 8px;
  padding: 12px;
  font-size: 14px;
  resize: none;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.05);
}

:deep(.el-button) {
  padding: 12px 24px;
  font-size: 14px;
  border-radius: 8px;
  background: linear-gradient(135deg, #1e3c72 0%, #2a5298 100%);
  border: none;
  transition: transform 0.2s;
}

:deep(.el-button:hover) {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(30, 60, 114, 0.2);
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@media screen and (max-width: 768px) {
  .chat-content {
    padding: 16px;
  }
  
  .chat-window {
    height: calc(100vh - 32px);
  }
  
  .message-list {
    padding: 16px;
  }
  
  .input-area {
    padding: 16px;
  }
}
</style> 