<template>
  <div>
    <!-- 页面头部 -->
    <div class="page-header">
      <h2 class="page-title">仪表板</h2>
      <p class="page-description">Docker镜像代理管理系统概览</p>
    </div>

    <!-- 导航卡片 -->
    <el-row :gutter="24">
      <el-col :xs="24" :sm="12" :lg="6">
        <el-card class="nav-card" shadow="hover" @click="$router.push('/registries')">
          <div class="nav-card-content">
            <div class="nav-card-icon registries">
              <el-icon size="32"><Connection /></el-icon>
            </div>
            <div class="nav-card-info">
              <h3>镜像源管理</h3>
              <p>管理上游Docker镜像源</p>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :xs="24" :sm="12" :lg="6">
        <el-card class="nav-card" shadow="hover" @click="$router.push('/whitelists')">
          <div class="nav-card-content">
            <div class="nav-card-icon whitelists">
              <el-icon size="32"><Key /></el-icon>
            </div>
            <div class="nav-card-info">
              <h3>白名单管理</h3>
              <p>配置镜像访问白名单</p>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :xs="24" :sm="12" :lg="6">
        <el-card class="nav-card" shadow="hover" @click="$router.push('/users')">
          <div class="nav-card-content">
            <div class="nav-card-icon users">
              <el-icon size="32"><UserFilled /></el-icon>
            </div>
            <div class="nav-card-info">
              <h3>用户管理</h3>
              <p>管理系统用户账户</p>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :xs="24" :sm="12" :lg="6">
        <el-card class="nav-card" shadow="hover" @click="$router.push('/logs')">
          <div class="nav-card-content">
            <div class="nav-card-icon logs">
              <el-icon size="32"><Document /></el-icon>
            </div>
            <div class="nav-card-info">
              <h3>访问日志</h3>
              <p>查看系统访问记录</p>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 统计信息 -->
    <div class="content-card" style="margin-top: 24px;">
      <div class="card-header">
        <h3 class="card-title">系统统计</h3>
      </div>
      <div class="card-body">
        <el-row :gutter="24">
          <el-col :xs="24" :sm="12" :lg="6">
            <div class="stat-item">
              <div class="stat-value">{{ stats.registries }}</div>
              <div class="stat-label">镜像源数量</div>
            </div>
          </el-col>
          <el-col :xs="24" :sm="12" :lg="6">
            <div class="stat-item">
              <div class="stat-value">{{ stats.whitelists }}</div>
              <div class="stat-label">白名单条目</div>
            </div>
          </el-col>
          <el-col :xs="24" :sm="12" :lg="6">
            <div class="stat-item">
              <div class="stat-value">{{ stats.users }}</div>
              <div class="stat-label">用户数量</div>
            </div>
          </el-col>
          <el-col :xs="24" :sm="12" :lg="6">
            <div class="stat-item">
              <div class="stat-value">{{ stats.logs }}</div>
              <div class="stat-label">今日访问</div>
            </div>
          </el-col>
        </el-row>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Connection, Key, UserFilled, Document } from '@element-plus/icons-vue'
import axios from 'axios'

const stats = ref({
  registries: 0,
  whitelists: 0,
  users: 0,
  logs: 0
})

const loadStats = async () => {
  try {
    const [registriesRes, whitelistsRes, usersRes, logsRes] = await Promise.all([
      axios.get('/api/registries'),
      axios.get('/api/whitelists'),
      axios.get('/api/users'),
      axios.get('/api/logs?limit=1000')
    ])
    
    stats.value.registries = registriesRes.data?.length || 0
    stats.value.whitelists = whitelistsRes.data?.length || 0
    stats.value.users = usersRes.data?.length || 0
    
    // 计算今日访问量
    const today = new Date().toDateString()
    const todayLogs = logsRes.data?.filter(log => {
      const logDate = new Date(log.created_at).toDateString()
      return logDate === today
    }) || []
    stats.value.logs = todayLogs.length
    
  } catch (error) {
    console.error('加载统计数据失败:', error)
  }
}

onMounted(() => {
  loadStats()
})
</script>

<style scoped>
.nav-card {
  margin-bottom: 16px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.nav-card:hover {
  transform: translateY(-2px);
}

.nav-card-content {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 8px;
}

.nav-card-icon {
  width: 64px;
  height: 64px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.nav-card-icon.registries {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.nav-card-icon.whitelists {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.nav-card-icon.users {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
}

.nav-card-icon.logs {
  background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
}

.nav-card-info h3 {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
  margin: 0 0 4px 0;
}

.nav-card-info p {
  font-size: 14px;
  color: #606266;
  margin: 0;
}

.stat-item {
  text-align: center;
  padding: 16px;
}

.stat-value {
  font-size: 28px;
  font-weight: 600;
  color: #409eff;
  margin-bottom: 8px;
}

.stat-label {
  font-size: 14px;
  color: #606266;
}

@media (max-width: 768px) {
  .nav-card-content {
    flex-direction: column;
    text-align: center;
    gap: 12px;
  }
  
  .nav-card-icon {
    width: 48px;
    height: 48px;
  }
  
  .stat-value {
    font-size: 24px;
  }
}
</style>
