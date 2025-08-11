<template>
  <div>
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <div class="header-main">
          <h2 class="page-title">仪表板</h2>
          <p class="page-description">Docker镜像代理系统</p>
        </div>
        <div class="version-info" v-if="version.version">
          <div class="project-info">
            <div class="version-badge">
              <span class="version-label">Beta</span>
              <span class="version-number">{{ version.version }}</span>
            </div>
            <div class="github-link">
              <el-button 
                text 
                size="small" 
                title="帮忙点个星星呗"
                @click="openGithub"
                class="github-button"
              >
                <el-icon><Star /></el-icon>
                GitHub
              </el-button>
            </div>
          </div>
          <!-- <div class="build-date">{{ formatDate(version.date) }}</div> -->
        </div>
      </div>
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
              <p>上游Docker镜像源</p>
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
              <p>管理系统用户</p>
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
import { Connection, Key, UserFilled, Document, Star } from '@element-plus/icons-vue'
import axios from 'axios'

const stats = ref({
  registries: 0,
  whitelists: 0,
  users: 0,
  logs: 0
})

const version = ref({
  version: '',
  date: ''
})

const openGithub = () => {
  window.open('https://github.com/helloxz/zmirror', '_blank')
}

const formatDate = (dateStr) => {
  if (!dateStr || dateStr.length !== 8) return ''
  const year = dateStr.substring(0, 4)
  const month = dateStr.substring(4, 6)
  const day = dateStr.substring(6, 8)
  return `${year}-${month}-${day}`
}

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

const loadVersion = async () => {
  try {
    const response = await axios.get('/api/version')
    version.value = response.data
  } catch (error) {
    console.error('加载版本信息失败:', error)
  }
}

onMounted(() => {
  loadStats()
  loadVersion()
})
</script>

<style scoped>
.page-header {
  margin-bottom: 24px;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 24px;
}

.header-main {
  flex: 1;
}

.version-info {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 8px;
}

.project-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.version-badge {
  display: flex;
  align-items: center;
  gap: 8px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 6px 12px;
  border-radius: 16px;
  font-size: 12px;
  font-weight: 500;
}

.github-link {
  display: flex;
  align-items: center;
}

.github-button {
  color: white !important;
  font-size: 12px !important;
  padding: 6px 12px !important;
  border-radius: 16px !important;
  background: linear-gradient(135deg, #24292e 0%, #4a5568 100%) !important;
  border: 1px solid rgba(255, 255, 255, 0.1) !important;
  transition: all 0.3s ease !important;
  font-weight: 500 !important;
}

.github-button:hover {
  background: linear-gradient(135deg, #2c3e50 0%, #34495e 100%) !important;
  transform: translateY(-1px) !important;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15) !important;
  border-color: rgba(255, 255, 255, 0.2) !important;
}

.github-button .el-icon {
  margin-right: 4px;
  font-size: 14px;
}

.version-label {
  opacity: 0.9;
}

.version-number {
  font-weight: 600;
}

.build-date {
  font-size: 11px;
  color: #909399;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
}

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
  .header-content {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
  }
  
  .version-info {
    align-items: flex-start;
  }
  
  .project-info {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }
  
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
