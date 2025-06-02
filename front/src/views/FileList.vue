<template>
  <div class="file-list-container">
    <div class="navbar">
      <div class="navbar-left">
        <img src="../assets/logo.svg" alt="Logo" class="navbar-logo" />
        <span class="navbar-title">Cloudisk 网盘系统</span>
      </div>
      <div class="navbar-right">
        <div class="user-info">
          <el-avatar :size="32" :icon="UserFilled" class="user-avatar" />
          <span class="navbar-user">{{ username }}</span>
        </div>
        <el-button type="primary" class="navbar-logout" @click="handleLogout">
          <el-icon><SwitchButton /></el-icon>
          <span>退出登录</span>
        </el-button>
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
            <span>上传管理</span>
          </div>
          <div class="sidebar-menu-item" :class="{active: currentView==='download'}" @click="currentView='download'">
            <el-icon><Download /></el-icon>
            <span>下载管理</span>
          </div>
          <div class="sidebar-menu-item" :class="{active: currentView==='ai-chat'}" @click="currentView='ai-chat'">
            <el-icon><ChatDotRound /></el-icon>
            <span>文件助手</span>
          </div>
          <div class="sidebar-menu-item" :class="{active: currentView==='recycle'}" @click="currentView='recycle'">
            <el-icon><Delete /></el-icon>
            <span>回收站</span>
          </div>
          <div class="sidebar-menu-item" :class="{active: currentView==='setting'}" @click="currentView='setting'">
            <el-icon><Setting /></el-icon>
            <span>设置</span>
          </div>
        </div>
      </div>
      <div class="main-content">
        <div v-if="currentView==='files'" class="file-list-content">
          <!-- 新建文件夹按钮和弹窗 -->
          <div class="filelist-toolbar">
            <el-button type="primary" @click="showMkdir = true"><el-icon><Folder /></el-icon> 新建文件夹</el-button>
            <el-button type="primary" class="navbar-upload" @click="showUpload = true"><el-icon><UploadFilled /></el-icon> 上传文件</el-button>
          </div>
          <el-dialog v-model="showMkdir" title="新建文件夹" width="360px" :close-on-click-modal="false">
            <el-form :model="mkdirForm" :rules="mkdirRules" ref="mkdirFormRef" label-width="80px">
              <el-form-item label="文件夹名" prop="dir_name">
                <el-input v-model="mkdirForm.dir_name" maxlength="100" show-word-limit placeholder="请输入文件夹名称" />
              </el-form-item>
            </el-form>
            <template #footer>
              <el-button @click="showMkdir = false">取消</el-button>
              <el-button type="primary" @click="handleMkdir">确定</el-button>
            </template>
          </el-dialog>
          <!-- 上传弹窗（仅选择文件） -->
          <el-dialog v-model="showUpload" title="上传文件" width="420px" :close-on-click-modal="true">
            <el-upload
              class="upload-demo"
              drag
              :auto-upload="true"
              :show-file-list="false"
              :http-request="customUploadRequest"
              :disabled="uploading"
            >
              <el-icon class="el-icon--upload"><upload-filled /></el-icon>
              <div class="el-upload__text">将文件拖到此处，或 <em>点击上传</em></div>
            </el-upload>
            <template #footer>
              <el-button @click="showUpload = false">关闭</el-button>
            </template>
          </el-dialog>
          <!-- 面包屑导航 -->
          <div class="breadcrumb">
            <el-breadcrumb separator="/">
              <el-breadcrumb-item :to="{ path: '/files', query: { father_id: '-1' } }">根目录</el-breadcrumb-item>
              <el-breadcrumb-item v-for="(folder, index) in folderPath" :key="index" 
                :to="{ path: '/files', query: { father_id: folder.user_file_id } }">
                {{ folder.file_name }}
              </el-breadcrumb-item>
            </el-breadcrumb>
          </div>
          <!-- 文件列表 -->
          <div class="file-list">
            <el-table :data="fileList" style="width: 100%" @row-click="handleRowClick">
              <el-table-column prop="file_name" label="文件名" min-width="200">
                <template #default="{ row }">
                  <div class="file-name">
                    <el-icon v-if="row.type === 2"><Folder /></el-icon>
                    <el-icon v-else><Document /></el-icon>
                    <span>{{ row.file_name }}</span>
                  </div>
                </template>
              </el-table-column>
              <el-table-column prop="file_size" label="大小" width="120">
                <template #default="{ row }">
                  {{ formatFileSize(row.file_size) }}
                </template>
              </el-table-column>
              <el-table-column prop="update_at" label="修改时间" width="180">
                <template #default="{ row }">
                  {{ formatDate(row.update_at) }}
                </template>
              </el-table-column>
              <el-table-column label="操作" width="150" fixed="right">
                <template #default="{ row }">
                  <el-button-group>
                    <el-button type="primary" link @click.stop="handleDownload(row)" v-if="row.type === 1">
                      下载
                    </el-button>
                    <el-button type="danger" link @click.stop="handleDelete(row)">
                      删除
                    </el-button>
                  </el-button-group>
                </template>
              </el-table-column>
            </el-table>
          </div>
        </div>
        <div v-else-if="currentView==='transfer'" class="transfer-content">
          <div class="transfer-title"><el-icon><UploadFilled /></el-icon> 传输管理</div>
          <div v-if="uploadTasks.length === 0" class="transfer-empty">暂无上传任务</div>
          <div v-else class="transfer-list">
            <div v-for="task in uploadTasks" :key="task.id" class="transfer-task">
              <div class="transfer-task-header">
                <span class="transfer-task-name">{{ task.fileName }}</span>
                <span class="transfer-task-status" :class="task.status">{{ transferStatusText(task) }}</span>
              </div>
              <el-progress :percentage="task.percent" :status="task.status === 'error' ? 'exception' : (task.status === 'success' ? 'success' : undefined)" />
              <div class="transfer-task-step">{{ task.stepText }}</div>
              <div class="transfer-task-actions">
                <el-button v-if="task.status === 'error'" size="small" @click="retryUpload(task)">重试</el-button>
                <el-button v-if="task.status === 'uploading'" size="small" @click="pauseUpload(task)" type="warning">暂停</el-button>
                <el-button v-if="task.status === 'paused'" size="small" @click="resumeUpload(task)" type="primary">继续</el-button>
                <el-button v-if="task.status === 'uploading'" size="small" @click="cancelUpload(task)" type="danger">取消</el-button>
                <el-button v-if="task.status === 'success'" size="small" type="danger" @click="removeUploadTask(task)"><el-icon><Delete /></el-icon> 删除</el-button>
              </div>
            </div>
          </div>
        </div>
        <div v-else-if="currentView==='download'" class="transfer-content">
          <div class="transfer-title"><el-icon><Download /></el-icon> 下载管理</div>
          <div v-if="downloadTasks.length === 0" class="transfer-empty">暂无下载任务</div>
          <div v-else class="transfer-list">
            <div v-for="task in downloadTasks" :key="task.id" class="transfer-task">
              <div class="transfer-task-header">
                <span class="transfer-task-name">{{ task.fileName }}</span>
                <span class="transfer-task-status" :class="task.status">{{ downloadStatusText(task) }}</span>
              </div>
              <el-progress :percentage="task.percent" :status="task.status === 'error' ? 'exception' : (task.status === 'success' ? 'success' : undefined)" />
              <div class="transfer-task-step">{{ task.stepText }}</div>
              <div class="transfer-task-actions">
                <el-button v-if="task.status === 'error'" size="small" @click="retryDownload(task)">重试</el-button>
                <el-button v-if="task.status === 'downloading'" size="small" @click="pauseDownload(task)" type="warning">暂停</el-button>
                <el-button v-if="task.status === 'paused'" size="small" @click="resumeDownload(task)" type="primary">继续</el-button>
                <el-button v-if="task.status === 'downloading'" size="small" @click="cancelDownload(task)" type="danger">取消</el-button>
                <el-button v-if="task.status === 'success'" size="small" type="danger" @click="removeDownloadTask(task)"><el-icon><Delete /></el-icon> 删除</el-button>
              </div>
            </div>
          </div>
        </div>
        <div v-else-if="currentView==='ai-chat'" class="chat-content">
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
              <!-- 加载状态 -->
              <div v-if="loading" class="message ai-message loading-message">
                <div class="message-avatar">
                  <el-avatar :size="40" :icon="Service" />
                </div>
                <div class="message-content">
                  <div class="message-text loading-text">
                    <span class="loading-dot"></span>
                    <span class="loading-dot"></span>
                    <span class="loading-dot"></span>
                  </div>
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
                :disabled="loading"
              />
              <div class="input-actions">
                <span class="tip">按 Ctrl + Enter 发送</span>
                <el-button 
                  type="primary" 
                  :loading="loading"
                  @click="sendMessage"
                  :disabled="loading"
                >
                  {{ loading ? '思考中...' : '发送' }}
                </el-button>
              </div>
            </div>
          </div>
        </div>
        <div v-else class="coming-soon">敬请期待...</div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Folder, Document, Delete, Setting, UploadFilled, Download, ChatDotRound, UserFilled, Service, SwitchButton } from '@element-plus/icons-vue'
import request from '../utils/request'
import SparkMD5 from 'spark-md5'
import axios from 'axios'
import { marked } from 'marked'
import DOMPurify from 'dompurify'

const route = useRoute()
const router = useRouter()
const fileList = ref([])
const folderPath = ref([])

// 用户名（假设token登录后用户名存在localStorage）
const username = ref(localStorage.getItem('username') || '用户')

// 上传相关
const showUpload = ref(false)
const uploading = ref(false)
const uploadPercent = ref(0)
const uploadStatus = ref('')
const uploadStatusText = ref('')
let currentFile = null
let canceling = false

const showTransfer = ref(false)
const uploadTasks = ref([])

// 上传任务结构: { id, file, fileName, status, percent, cancelToken, hash }
let uploadTaskId = 0

const currentView = ref('files')

const showMkdir = ref(false)
const mkdirForm = ref({ dir_name: '' })
const mkdirFormRef = ref()
const mkdirRules = {
  dir_name: [
    { required: true, message: '请输入文件夹名称', trigger: 'blur' },
    { min: 1, max: 100, message: '长度1-100字符', trigger: 'blur' }
  ]
}

// 下载相关
const downloadTasks = ref([])
let downloadTaskId = 0

const handleLogout = () => {
  localStorage.removeItem('token')
  router.push('/login')
}

// 获取文件夹路径（前端维护）
function updateFolderPathOnEnter(row) {
  // row为当前点击的文件夹对象
  const idx = folderPath.value.findIndex(f => f.user_file_id === row.user_file_id)
  if (idx !== -1) {
    // 已存在，截断到该层级
    folderPath.value = folderPath.value.slice(0, idx + 1)
  } else {
    // 追加
    folderPath.value.push({
      user_file_id: row.user_file_id,
      file_name: row.file_name
    })
  }
}

// 获取文件列表
const getFileList = async (fatherId) => {
  console.log('Attempting to fetch file list for father_id:', fatherId);
  console.trace('getFileList call stack'); // Keep the trace for now

  // *** 新增检查 ***
  if (fatherId === undefined || fatherId === null) {
    console.warn('getFileList called with undefined or null father_id. Request aborted.');
    return; // Abort the function if fatherId is invalid
  }
  // *** 结束新增检查 ***

  try {
    // 保持原有的逻辑，确保 '-1' 被正确处理
    const actualFatherId = (fatherId === '-1' || fatherId === -1) ? '-1' : fatherId;

    console.log('Fetching file list for actual father_id:', actualFatherId); // Log the actual ID used
    const response = await request.get('/file', {
      params: { father_id: actualFatherId }
    });
    fileList.value = response.data
    if (!actualFatherId || actualFatherId === '-1') {
      folderPath.value = []
    } else {
      // 如果不是根目录，且folderPath是空的，可能是在子目录刷新的情况，需要获取完整路径
      if(folderPath.value.length === 0) {
         // 这里如果需要完整路径，可能还是需要一个接口来获取，或者调整前端路径维护逻辑
         // 暂时不处理刷新子目录后路径不完整的问题，只保证请求参数正确
      }
    }

  } catch (error) {
    // 错误处理已经在请求拦截器中完成
  }
}

// 处理行点击
const handleRowClick = (row) => {
  if (row.type === 2) { // 如果是文件夹
    router.push({
      path: '/files',
      query: { father_id: row.user_file_id }
    })
    updateFolderPathOnEnter(row)
  }
}

// 下载文件
const downloadFile = async (task) => {
  try {
    const token = localStorage.getItem('token')
    task.stepText = '正在连接服务器...'
    console.log('开始下载文件:', task.fileName)
    
    let retryCount = 0
    const maxRetries = 3

    while (retryCount < maxRetries) {
      try {
        // 创建取消令牌
        const source = axios.CancelToken.source()
        task.cancelToken = source

        const response = await axios({
          url: '/file/download',
          method: 'get',
          params: { user_file_id: task.file.user_file_id },
          responseType: 'blob',
          headers: {
            'Accept': 'application/octet-stream',
            'X-Token': token
          },
          baseURL: '/api',
          cancelToken: source.token,
          timeout: 30000, // 30秒超时
          onDownloadProgress: (progressEvent) => {
            console.log('下载进度:', progressEvent)
            if (progressEvent.lengthComputable) {
              const percent = Math.round((progressEvent.loaded / progressEvent.total) * 100)
              task.percent = percent
              task.stepText = `已下载 ${formatFileSize(progressEvent.loaded)} / ${formatFileSize(progressEvent.total)} (${percent}%)`
              // 强制更新视图
              downloadTasks.value = [...downloadTasks.value]
            } else {
              // 如果无法计算总大小，至少显示已下载的大小
              task.stepText = `已下载 ${formatFileSize(progressEvent.loaded)}`
              // 强制更新视图
              downloadTasks.value = [...downloadTasks.value]
            }
          }
        })

        console.log('收到响应:', response.headers)

        // 检查响应类型
        const contentType = response.headers['content-type']
        if (contentType && contentType.includes('application/json')) {
          const reader = new FileReader()
          reader.onload = () => {
            try {
              const errorData = JSON.parse(reader.result)
              throw new Error(errorData.message || '下载失败')
            } catch (e) {
              throw new Error('下载失败')
            }
          }
          reader.readAsText(response.data)
          return
        }

        // 创建Blob对象
        const blob = new Blob([response.data], { type: contentType || 'application/octet-stream' })
        
        // 检查文件大小
        if (blob.size === 0) {
          throw new Error('下载的文件为空')
        }

        task.stepText = '正在保存文件...'
        // 强制更新视图
        downloadTasks.value = [...downloadTasks.value]

        // 创建下载链接
        const url = window.URL.createObjectURL(blob)
        const link = document.createElement('a')
        link.style.display = 'none'
        link.href = url
        link.download = task.fileName
        document.body.appendChild(link)
        link.click()
        
        // 清理
        setTimeout(() => {
          document.body.removeChild(link)
          window.URL.revokeObjectURL(url)
        }, 100)

        task.status = 'success'
        task.percent = 100
        task.stepText = '下载完成'
        // 强制更新视图
        downloadTasks.value = [...downloadTasks.value]
        break // 下载成功，跳出重试循环

      } catch (error) {
        console.error(`下载尝试 ${retryCount + 1} 失败:`, error)
        if (axios.isCancel(error)) {
          console.log('下载被取消')
          return
        }
        if (retryCount < maxRetries - 1) {
          retryCount++
          task.stepText = `连接失败，正在重试(${retryCount}/${maxRetries})...`
          // 强制更新视图
          downloadTasks.value = [...downloadTasks.value]
          await new Promise(resolve => setTimeout(resolve, 1000 * retryCount)) // 递增延迟
          continue
        }
        throw error
      }
    }
  } catch (error) {
    console.error('下载错误:', error)
    if (error.response) {
      console.error('错误响应:', error.response)
    }
    task.status = 'error'
    task.stepText = error.message || '下载失败'
    // 强制更新视图
    downloadTasks.value = [...downloadTasks.value]
  }
}

// 处理下载
const handleDownload = (file) => {
  // 获取token
  const token = localStorage.getItem('token')
  if (!token) {
    ElMessage.error('请先登录')
    router.push('/login')
    return
  }

  // 创建下载任务
  const task = {
    id: ++downloadTaskId,
    file: file,
    fileName: file.file_name,
    status: 'downloading',
    percent: 0,
    stepText: '准备下载...',
    cancelToken: { canceled: false, paused: false }
  }
  downloadTasks.value.push(task)

  // 切换到下载管理页面
  currentView.value = 'download'

  // 开始下载
  downloadFile(task)
}

// 下载状态文本
function downloadStatusText(task) {
  if (task.status === 'downloading') return '下载中'
  if (task.status === 'success') return '完成'
  if (task.status === 'error') return '失败'
  if (task.status === 'paused') return '已暂停'
  return task.status
}

// 重试下载
function retryDownload(task) {
  if (task.status === 'error') {
    task.status = 'downloading'
    task.percent = 0
    task.cancelToken = { canceled: false, paused: false }
    downloadFile(task)
  }
}

// 暂停下载
function pauseDownload(task) {
  if (task.status === 'downloading') {
    task.status = 'paused'
    task.cancelToken.paused = true
    task.stepText = '已暂停'
  }
}

// 继续下载
function resumeDownload(task) {
  if (task.status === 'paused') {
    task.status = 'downloading'
    task.cancelToken.paused = false
    task.stepText = '继续下载...'
    downloadFile(task)
  }
}

// 取消下载
function cancelDownload(task) {
  if (task.status === 'downloading') {
    if (task.cancelToken) {
      task.cancelToken.cancel('用户取消下载')
    }
    task.status = 'error'
    task.percent = 0
    task.stepText = '已取消'
    // 强制更新视图
    downloadTasks.value = [...downloadTasks.value]
  }
}

// 移除下载任务
function removeDownloadTask(task) {
  const idx = downloadTasks.value.findIndex(t => t.id === task.id)
  if (idx !== -1) downloadTasks.value.splice(idx, 1)
}

// 处理删除
const handleDelete = async (file) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除${file.type === 2 ? '文件夹' : '文件'} "${file.file_name}" 吗？`,
      '警告',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    // 在异步操作前获取当前的father_id
    const currentFatherId = route.query.father_id;
    const response = await request.post('/file/remove', {
      user_file_id: file.user_file_id,
      type: file.type
    })
    if (response.code === 200) {
      ElMessage.success('删除成功')
      // 删除成功后，先从当前列表中移除该项，提高用户体验
      const index = fileList.value.findIndex(item => item.user_file_id === file.user_file_id)
      if (index !== -1) {
        fileList.value.splice(index, 1)
      }
      // 然后再调用接口刷新列表，使用之前获取的father_id
      console.log('Deleting file/folder successful, refetching file list for father_id:', currentFatherId);
      getFileList(currentFatherId)
    }
  } catch (error) {
    // 错误提示已由拦截器处理
  }
}

// 格式化文件大小
const formatFileSize = (size) => {
  if (size === 0) return '-'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  let index = 0
  while (size >= 1024 && index < units.length - 1) {
    size /= 1024
    index++
  }
  return `${size.toFixed(2)} ${units[index]}`
}

// 格式化日期
const formatDate = (dateStr) => {
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 状态文本
function transferStatusText(task) {
  if (task.status === 'waiting') return '等待中'
  if (task.status === 'uploading') return '上传中'
  if (task.status === 'success') return '完成'
  if (task.status === 'error') return '失败'
  return task.status
}

// 上传弹窗自定义上传
const customUploadRequest = async (option) => {
  uploadTasks.value.push({
    id: ++uploadTaskId,
    file: option.file,
    fileName: option.file.name,
    status: 'waiting',
    percent: 0,
    hash: '',
    cancelToken: { canceled: false, paused: false },
    stepText: ''
  })
  showUpload.value = false
  ElMessage.info('文件已加入上传队列，可在传输管理页面查看进度')
  startNextUpload()
}

// 启动下一个等待中的任务（支持并发可扩展）
function startNextUpload() {
  const task = uploadTasks.value.find(t => t.status === 'waiting')
  if (task) {
    uploadFileTask(task)
  }
}

// 上传单个任务
async function uploadFileTask(task) {
  task.status = 'uploading'
  task.percent = 0
  try {
    // 1. 计算hash
    task.stepText = '正在计算文件Hash...'
    task.hash = await calcFileHash(task.file)
    if (task.cancelToken.canceled) throw new Error('已取消')
    if (task.cancelToken.paused) {
      task.status = 'paused'
      task.stepText = '已暂停'
      return
    }
    // 2. 初始化
    task.stepText = '正在初始化上传...'
    // 获取并转换为数字类型的father_id
    const fatherIdInt = Number(route.query.father_id || -1);
    const initRes = await request.post('/file/init', {
      file_name: task.file.name,
      file_size: task.file.size,
      file_hash: task.hash,
      father_id: fatherIdInt // 使用数字类型的father_id
    })
    if (initRes.code === 200 && initRes.message && initRes.message.includes('秒传')) {
      task.percent = 100
      task.status = 'success'
      task.stepText = '秒传成功！'
      getFileList(fatherIdInt)
      startNextUpload()
      return
    }
    if (!initRes.data || typeof initRes.data.chunk_size !== 'number' || typeof initRes.data.chunk_num !== 'number') {
      throw new Error('初始化失败，参数不完整')
    }
    // 3. 分片上传
    const { chunk_size, chunk_num, has_chunk = [] } = initRes.data
    for (let i = 0; i < chunk_num; i++) {
      if (task.cancelToken.canceled) throw new Error('已取消')
      if (task.cancelToken.paused) {
        task.status = 'paused'
        task.stepText = '已暂停'
        return
      }
      if (has_chunk && has_chunk.includes(String(i))) {
        task.percent = Math.round(((i + 1) / chunk_num) * 100)
        continue
      }
      const chunk = task.file.slice(i * chunk_size, (i + 1) * chunk_size)
      const formData = new FormData()
      formData.append('chunk_idx', String(i))
      formData.append('file_hash', task.hash)
      formData.append('file', chunk)
      let retry = 0
      task.stepText = `正在上传分片 ${i + 1} / ${chunk_num}`
      while (retry < 3) {
        try {
          const res = await request.post('/file/chunk', formData)
          if (res.code === 200) break
          else throw new Error(res.message || '分片上传失败')
        } catch (err) {
          retry++
          if (retry >= 3) throw err
          await new Promise(r => setTimeout(r, 500 * retry))
        }
      }
      task.percent = Math.round(((i + 1) / chunk_num) * 100)
    }
    // 4. 合并分片
    task.stepText = '正在合并分片...'
    const mergeRes = await request.post('/file/merge', { file_hash: task.hash })
    if (mergeRes.code === 200) {
      task.percent = 100
      task.status = 'success'
      task.stepText = '上传完成！'
      getFileList(fatherIdInt)
      startNextUpload()
    } else {
      throw new Error(mergeRes.message || '合并分片失败')
    }
  } catch (e) {
    if (task.cancelToken.canceled) {
      task.status = 'canceled'
      task.percent = 0
      task.stepText = '已取消'
    } else if (task.status === 'paused') {
      // 保持暂停状态
    } else {
      task.status = 'error'
      task.percent = 0
      task.stepText = '上传失败：' + (e?.message || e)
    }
    startNextUpload()
  }
}

// 重试上传
function retryUpload(task) {
  if (task.status === 'error') {
    task.status = 'waiting'
    task.percent = 0
    task.cancelToken = { canceled: false }
    startNextUpload()
  }
}

// 暂停上传
function pauseUpload(task) {
  if (task.status === 'uploading') {
    task.status = 'paused'
    task.cancelToken.paused = true
    task.stepText = '已暂停'
  }
}

// 继续上传
function resumeUpload(task) {
  if (task.status === 'paused') {
    task.status = 'waiting'
    task.cancelToken.paused = false
    task.stepText = '等待中...'
    startNextUpload()
  }
}

// 取消上传
function cancelUpload(task) {
  if (task && task.status === 'uploading') {
    task.cancelToken.canceled = true
    task.status = 'error'
    task.percent = 0
  }
}

// 计算文件hash（MD5）
function calcFileHash(file) {
  return new Promise((resolve, reject) => {
    const chunkSize = 2 * 1024 * 1024
    const chunks = Math.ceil(file.size / chunkSize)
    let currentChunk = 0
    const spark = new SparkMD5.ArrayBuffer()
    const fileReader = new FileReader()
    fileReader.onload = function (e) {
      spark.append(e.target.result)
      currentChunk++
      if (currentChunk < chunks) {
        loadNext()
      } else {
        resolve(spark.end())
      }
    }
    fileReader.onerror = function () {
      reject('文件读取失败')
    }
    function loadNext() {
      if (canceling) return
      const start = currentChunk * chunkSize
      const end = Math.min(file.size, start + chunkSize)
      fileReader.readAsArrayBuffer(file.slice(start, end))
    }
    loadNext()
  })
}

// 监听路由参数变化
watch(
  () => route.query.father_id,
  (newFatherId) => {
    getFileList(newFatherId)
    if (!newFatherId || newFatherId === '-1') {
      folderPath.value = []
    }
  }
)

// 组件挂载时获取文件列表
onMounted(() => {
  const fatherId = route.query.father_id || '-1'
  getFileList(fatherId)
})

function removeUploadTask(task) {
  const idx = uploadTasks.value.findIndex(t => t.id === task.id)
  if (idx !== -1) uploadTasks.value.splice(idx, 1)
}

async function handleMkdir() {
  await mkdirFormRef.value.validate()
  const father_id = route.query.father_id || -1
  try {
    const res = await request.post('/file/mkdir', {
      father_id: Number(father_id),
      dir_name: mkdirForm.value.dir_name.trim()
    })
    if (res.code === 200) {
      ElMessage.success('创建文件夹成功')
      showMkdir.value = false
      mkdirForm.value.dir_name = ''
      getFileList(father_id)
    }
  } catch (e) {
    // 错误提示已由拦截器处理
  }
}

// AI 聊天相关
const messages = ref([])
const inputMessage = ref('')
const loading = ref(false)
const messageList = ref(null)

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
    console.log('正在发送消息:', message) // 添加日志
    const response = await request.post('/ai/chat', {
      message: message
    })
    console.log('收到响应:', response) // 添加日志
    
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
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 32px;
  z-index: 10;
  position: relative;
  backdrop-filter: blur(10px);
  background: rgba(255, 255, 255, 0.95);
}

.navbar-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.navbar-logo {
  height: 36px;
  transition: transform 0.3s ease;
}

.navbar-logo:hover {
  transform: scale(1.05);
}

.navbar-title {
  font-size: 20px;
  font-weight: 600;
  background: linear-gradient(135deg, #1e3c72 0%, #2a5298 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  letter-spacing: 0.5px;
}

.navbar-right {
  display: flex;
  align-items: center;
  gap: 24px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 4px 12px;
  border-radius: 20px;
  background: #f5f7fa;
  transition: all 0.3s ease;
}

.user-info:hover {
  background: #e4e7eb;
}

.user-avatar {
  background: linear-gradient(135deg, #1e3c72 0%, #2a5298 100%);
  color: #fff;
}

.navbar-user {
  font-size: 14px;
  color: #1e3c72;
  font-weight: 500;
}

.navbar-logout {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  background: linear-gradient(135deg, #1e3c72 0%, #2a5298 100%);
  border: none;
  transition: all 0.3s ease;
}

.navbar-logout:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(30, 60, 114, 0.2);
}

.navbar-logout .el-icon {
  font-size: 16px;
}

@media screen and (max-width: 768px) {
  .navbar {
    padding: 0 16px;
  }
  
  .navbar-title {
    font-size: 18px;
  }
  
  .navbar-user {
    display: none;
  }
  
  .navbar-logout span {
    display: none;
  }
  
  .navbar-logout {
    padding: 8px;
  }
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

.file-list-content {
  width: 100%;
  max-width: 900px;
  background: #fff;
  border-radius: 16px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
  overflow: hidden;
  display: flex;
  flex-direction: column;
  gap: 20px;
  padding: 24px;
  margin-top: 32px;
}

.breadcrumb {
  padding: 16px 24px;
  background: #f8f9fa;
  border-radius: 12px;
  margin-bottom: 8px;
}

.breadcrumb :deep(.el-breadcrumb__item) {
  display: flex;
  align-items: center;
}

.breadcrumb :deep(.el-icon) {
  margin-right: 4px;
  font-size: 16px;
}

.breadcrumb :deep(.el-breadcrumb__inner) {
  font-size: 14px;
  color: #606266;
}

.breadcrumb :deep(.el-breadcrumb__inner.is-link:hover) {
  color: #409EFF;
}

.file-list {
  flex: 1;
  background: #fff;
  border-radius: 12px;
  overflow: hidden;
}

.file-list :deep(.el-table) {
  --el-table-border-color: #ebeef5;
  --el-table-header-bg-color: #f8f9fa;
}

.file-list :deep(.el-table th) {
  font-weight: 600;
  color: #606266;
  background: #f8f9fa;
  padding: 16px 0;
}

.file-list :deep(.el-table td) {
  padding: 16px 0;
}

.file-name {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 0 8px;
}

.file-name .el-icon {
  font-size: 20px;
}

.folder-icon {
  color: #409EFF;
}

.file-icon {
  color: #909399;
}

:deep(.el-table__row) {
  cursor: pointer;
  transition: all 0.3s ease;
}

:deep(.el-table__row:hover) {
  background-color: #f5f7fa;
  transform: translateY(-1px);
}

:deep(.el-button-group .el-button) {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 8px 12px;
  font-size: 13px;
}

:deep(.el-button-group .el-button:hover) {
  transform: translateY(-1px);
}

@media screen and (max-width: 768px) {
  .file-list-content {
    padding: 8px;
    max-width: 98vw;
  }
  .breadcrumb {
    padding: 8px 8px;
  }
}

.navbar-upload {
  margin-right: 16px;
}
.upload-progress {
  margin: 16px 0 0 0;
}
.upload-status-text {
  margin-top: 8px;
  color: #666;
  font-size: 14px;
  text-align: center;
}

.navbar-transfer {
  margin-right: 16px;
}
.transfer-empty {
  text-align: center;
  color: #aaa;
  padding: 40px 0;
  font-size: 16px;
}

.sidebar-transfer {
  width: 100%;
  margin-top: 32px;
  padding: 0 8px 16px 8px;
  border-top: 1px solid #f0f0f0;
}
.sidebar-transfer-title {
  font-size: 16px;
  font-weight: 600;
  color: #1e3c72;
  margin: 16px 0 12px 0;
  display: flex;
  align-items: center;
  gap: 6px;
}
.transfer-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}
.transfer-task {
  background: #f8f9fa;
  border-radius: 8px;
  padding: 12px 10px 8px 10px;
  box-shadow: 0 2px 8px 0 rgba(30,60,114,0.04);
}
.transfer-task-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 4px;
}
.transfer-task-name {
  font-size: 14px;
  font-weight: 500;
  color: #1e3c72;
}
.transfer-task-status {
  font-size: 13px;
  color: #888;
}
.transfer-task-status.success {
  color: #67c23a;
}
.transfer-task-status.error {
  color: #f56c6c;
}
.transfer-task-status.uploading {
  color: #409EFF;
}
.transfer-task-status.waiting {
  color: #aaa;
}
.transfer-task-step {
  font-size: 12px;
  color: #888;
  margin: 4px 0 0 0;
  min-height: 18px;
}
.transfer-task-actions {
  margin-top: 4px;
  display: flex;
  gap: 8px;
}

.transfer-content {
  width: 100%;
  max-width: 900px;
  background: #fff;
  border-radius: 16px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
  overflow: hidden;
  display: flex;
  flex-direction: column;
  gap: 20px;
  padding: 24px;
  margin-top: 32px;
}
.transfer-title {
  font-size: 20px;
  font-weight: 600;
  color: #1e3c72;
  margin-bottom: 16px;
  display: flex;
  align-items: center;
  gap: 8px;
}
.coming-soon {
  width: 100%;
  max-width: 900px;
  background: #fff;
  border-radius: 16px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 22px;
  color: #aaa;
  margin-top: 32px;
  height: 300px;
}

.filelist-toolbar {
  display: flex;
  align-items: center;
  margin-bottom: 16px;
  gap: 12px;
}

/* AI 聊天相关样式 */
.chat-content {
  width: 100%;
  max-width: 900px;
  padding: 32px;
}

.chat-window {
  width: 100%;
  height: calc(100vh - 128px);
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
    height: calc(100vh - 96px);
  }
  
  .message-list {
    padding: 16px;
  }
  
  .input-area {
    padding: 16px;
  }
}

/* 加载状态样式 */
.loading-message {
  opacity: 0.8;
}

.loading-text {
  display: flex;
  align-items: center;
  gap: 4px;
  min-width: 60px;
}

.loading-dot {
  width: 8px;
  height: 8px;
  background: #1e3c72;
  border-radius: 50%;
  animation: loadingDot 1.4s infinite ease-in-out both;
}

.loading-dot:nth-child(1) {
  animation-delay: -0.32s;
}

.loading-dot:nth-child(2) {
  animation-delay: -0.16s;
}

@keyframes loadingDot {
  0%, 80%, 100% { 
    transform: scale(0);
  } 
  40% { 
    transform: scale(1.0);
  }
}

/* 禁用状态样式 */
:deep(.el-textarea.is-disabled .el-textarea__inner) {
  background-color: #f5f7fa;
  border-color: #e4e7ed;
  color: #606266;
  cursor: not-allowed;
}

:deep(.el-button.is-disabled) {
  background: #a0cfff;
  border-color: #a0cfff;
  color: #fff;
  cursor: not-allowed;
}
</style> 