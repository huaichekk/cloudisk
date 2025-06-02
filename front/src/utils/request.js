import axios from 'axios'
import { ElMessage } from 'element-plus'
import router from '../router'

// 创建 axios 实例
const request = axios.create({
    baseURL: '/api', // 使用相对路径，会被代理到后端服务器
    timeout: 30000, // 修改为 30 秒
    headers: {
        'Content-Type': 'application/x-www-form-urlencoded',
        'Accept': 'application/json'
    }
})

// 请求拦截器
request.interceptors.request.use(
    config => {
        const token = localStorage.getItem('token')
        if (token) {
            config.headers['X-Token'] = token
        }

        // 如果是 AI 聊天接口，设置更长的超时时间
        if (config.url === '/ai/chat') {
            config.timeout = 60000 // AI 聊天接口使用 60 秒超时
        }

        // 如果是 FormData，删除 Content-Type，让浏览器自动设置
        if (config.data instanceof FormData) {
            delete config.headers['Content-Type']
        } else if (config.data && typeof config.data === 'object' && config.headers['Content-Type']) {
            // 如果是json对象，自动设置为application/json
            config.headers['Content-Type'] = 'application/json'
        }

        return config
    },
    error => {
        return Promise.reject(error)
    }
)

// 响应拦截器
request.interceptors.response.use(
    response => {
        const res = response.data
        if (res.code === 0 || res.code === 200) {
            return res
        } else {
            ElMessage.error(res.msg || '请求失败')
            return Promise.reject(new Error(res.msg || '请求失败'))
        }
    },
    error => {
        if (error.response) {
            switch (error.response.status) {
                case 401:
                    ElMessage.error('登录已过期，请重新登录')
                    localStorage.removeItem('token')
                    router.push('/login')
                    break
                case 403:
                    ElMessage.error('没有权限访问')
                    break
                case 404:
                    ElMessage.error('请求的资源不存在')
                    break
                case 500:
                    ElMessage.error('服务器错误')
                    break
                default:
                    ElMessage.error(error.response.data?.msg || '请求失败')
            }
        } else {
            ElMessage.error('网络错误，请检查您的网络连接')
        }
        return Promise.reject(error)
    }
)

export default request 