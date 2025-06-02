import './assets/main.css'
import { createApp } from 'vue'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import App from './App.vue'
import router from './router'
import axios from 'axios'

// 配置axios默认值
axios.defaults.baseURL = 'http://localhost:8080' // 根据实际后端地址修改

// 添加请求拦截器
axios.interceptors.request.use(
    config => {
        const token = localStorage.getItem('token')
        if (token) {
            config.headers['X-Token'] = token
        }
        return config
    },
    error => {
        return Promise.reject(error)
    }
)

const app = createApp(App)

app.use(ElementPlus)
app.use(router)

app.mount('#app')
