<template>
  <div class="layout-container">
    <!-- 主界面 -->
    <template v-if="authStore.isLoggedIn && $route.name !== 'Home'">
      <!-- 头部 -->
      <header class="main-header">
        <div class="header-content">
          <div>
            <h1 class="header-title">Zmirror</h1>
            <p class="header-subtitle">Docker镜像代理系统</p>
          </div>
          <div class="header-actions">
            <div class="user-info">
              <el-icon><User /></el-icon>
              <span>{{ authStore.username || '管理员' }}</span>
            </div>
            <el-button
              type="danger"
              size="small"
              plain
              @click="handleLogout"
            >
              <el-icon><SwitchButton /></el-icon>
              退出登录
            </el-button>
          </div>
        </div>
      </header>

      <!-- 主体内容 -->
      <main class="main-content">
        <div class="content-container">
          <!-- 导航菜单 -->
          <el-menu
            :default-active="$route.path === '/' ? '/dashboard' : $route.path"
            mode="horizontal"
            router
            style="margin-bottom: 24px; border-radius: 8px; background: white; padding: 0 16px;"
          >
            <el-menu-item index="/dashboard">
              <el-icon><Odometer /></el-icon>
              仪表板
            </el-menu-item>
            <el-menu-item index="/registries">
              <el-icon><Connection /></el-icon>
              镜像源管理
            </el-menu-item>
            <el-menu-item index="/whitelists">
              <el-icon><Key /></el-icon>
              白名单管理
            </el-menu-item>
            <el-menu-item index="/users">
              <el-icon><UserFilled /></el-icon>
              用户管理
            </el-menu-item>
            <el-menu-item index="/logs">
              <el-icon><Document /></el-icon>
              访问日志
            </el-menu-item>
          </el-menu>
          
          <router-view />
        </div>
      </main>
    </template>
    
    <!-- 未登录时或首页显示路由内容 -->
    <template v-else>
      <router-view />
    </template>
  </div>
</template>

<script setup>
import { ElMessage } from 'element-plus'
import { useRouter } from 'vue-router'
import { useAuthStore } from './stores/auth'
import { 
  User, 
  SwitchButton, 
  Odometer, 
  Connection, 
  Key, 
  UserFilled, 
  Document 
} from '@element-plus/icons-vue'

const router = useRouter()
const authStore = useAuthStore()

// 退出登录
const handleLogout = () => {
  authStore.logout()
  ElMessage.success('已退出登录')
  router.push('/login')
}
</script>

<style scoped>
.layout-container {
  min-height: 100vh;
  background: #f5f7fa;
}

.main-header {
  background: white;
  border-bottom: 1px solid #e4e7ed;
  position: sticky;
  top: 0;
  z-index: 1000;
}

.header-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 16px 24px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.header-title {
  font-size: 24px;
  font-weight: 600;
  color: #303133;
  margin: 0 0 4px 0;
}

.header-subtitle {
  font-size: 14px;
  color: #606266;
  margin: 0;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 16px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: #606266;
}

.main-content {
  max-width: 1200px;
  margin: 0 auto;
  /* padding: 16px; */
}

.content-container {
  width:1200px;

}

@media (max-width: 768px) {
  .header-content {
    padding: 12px 16px;
    flex-direction: column;
    gap: 12px;
    text-align: center;
  }
  
  .main-content {
    padding: 16px;
  }
  
  .content-container {
    min-width: auto; /* 移动端取消最小宽度限制 */
  }
  
  .header-title {
    font-size: 20px;
  }
}
</style>
