import { defineStore } from 'pinia'
import { ref } from 'vue'
import axios from 'axios'

export const useAuthStore = defineStore('auth', () => {
  const isLoggedIn = ref(false)
  const username = ref('')
  const token = ref('')

  // 初始化
  const init = () => {
    const savedToken = localStorage.getItem('auth_token')
    const savedUsername = localStorage.getItem('auth_username')
    
    if (savedToken && savedUsername) {
      token.value = savedToken
      username.value = savedUsername
      isLoggedIn.value = true
      
      // 设置axios默认header
      axios.defaults.headers.common['Authorization'] = savedToken
    }
  }

  // 登录
  const login = async (user, password) => {
    const authToken = 'Basic ' + btoa(user + ':' + password)
    
    try {
      // 测试认证
      await axios.get('/api/registries', {
        headers: {
          'Authorization': authToken
        }
      })
      
      // 登录成功
      token.value = authToken
      username.value = user
      isLoggedIn.value = true
      
      // 保存到localStorage
      localStorage.setItem('auth_token', authToken)
      localStorage.setItem('auth_username', user)
      
      // 设置axios默认header
      axios.defaults.headers.common['Authorization'] = authToken
      
    } catch (error) {
      throw new Error('用户名或密码错误')
    }
  }

  // 退出登录
  const logout = () => {
    isLoggedIn.value = false
    username.value = ''
    token.value = ''
    
    // 清除localStorage
    localStorage.removeItem('auth_token')
    localStorage.removeItem('auth_username')
    
    // 清除axios默认header
    delete axios.defaults.headers.common['Authorization']
  }

  // 初始化
  init()

  return {
    isLoggedIn,
    username,
    token,
    login,
    logout
  }
})
