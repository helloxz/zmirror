<template>
  <div>
    <!-- 页面头部 -->
    <div class="page-header">
      <h2 class="page-title">访问日志</h2>
      <p class="page-description">查看系统访问记录和用户行为</p>
    </div>

    <!-- 内容卡片 -->
    <div class="content-card">
      <div class="card-header">
        <h3 class="card-title">访问记录</h3>
        <div style="display: flex; gap: 12px;">
          <el-button @click="loadData">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
          <el-button @click="handleClearLogs" type="danger" plain>
            <el-icon><Delete /></el-icon>
            清空日志
          </el-button>
        </div>
      </div>
      
      <div class="card-body">
        <!-- 筛选器 -->
        <div class="filter-container">
          <el-row :gutter="16">
            <el-col :span="6">
              <el-select
                v-model="filters.method"
                placeholder="请求方法"
                clearable
                style="width: 100%"
                @change="handleFilterChange"
              >
                <el-option label="GET" value="GET" />
                <el-option label="POST" value="POST" />
                <el-option label="PUT" value="PUT" />
                <el-option label="DELETE" value="DELETE" />
                <el-option label="HEAD" value="HEAD" />
              </el-select>
            </el-col>
            <el-col :span="6">
              <el-select
                v-model="filters.status"
                placeholder="状态码"
                clearable
                style="width: 100%"
                @change="handleFilterChange"
              >
                <el-option label="200 - 成功" value="200" />
                <el-option label="401 - 未授权" value="401" />
                <el-option label="403 - 禁止访问" value="403" />
                <el-option label="404 - 未找到" value="404" />
                <el-option label="500 - 服务器错误" value="500" />
              </el-select>
            </el-col>
            <el-col :span="6">
              <el-input
                v-model="filters.username"
                placeholder="用户名"
                clearable
                @keyup.enter="handleFilterChange"
              />
            </el-col>
            <el-col :span="6">
              <el-input
                v-model="filters.clientIp"
                placeholder="客户端IP"
                clearable
                @keyup.enter="handleFilterChange"
              />
            </el-col>
          </el-row>
        </div>

        <!-- 表格 -->
        <div class="table-container">
          <el-table
            v-loading="loading"
            :data="filteredLogs"
            stripe
            style="width: 100%"
            max-height="600"
            empty-text="暂无访问日志"
          >
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="client_ip" label="客户端IP" width="140" align="center">
              <template #default="{ row }">
                <el-tag size="small" type="info">{{ row.client_ip }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="method" label="方法" width="80" align="center">
              <template #default="{ row }">
                <el-tag
                  size="small"
                  :type="getMethodTagType(row.method)"
                >
                  {{ row.method }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="path" label="路径" min-width="300">
              <template #default="{ row }">
                <div style="font-family: monospace; font-size: 12px; word-break: break-all;">
                  {{ row.path }}
                </div>
              </template>
            </el-table-column>
            <el-table-column prop="status_code" label="状态码" width="100" align="center">
              <template #default="{ row }">
                <el-tag
                  size="small"
                  :type="getStatusTagType(row.status_code)"
                >
                  {{ row.status_code }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="username" label="用户" width="120" align="center">
              <template #default="{ row }">
                <span v-if="row.username">
                  <el-icon><User /></el-icon>
                  {{ row.username }}
                </span>
                <span v-else style="color: #909399;">匿名</span>
              </template>
            </el-table-column>
            <el-table-column prop="created_at" label="时间" width="180" align="center">
              <template #default="{ row }">
                {{ formatDate(row.created_at) }}
              </template>
            </el-table-column>
          </el-table>
        </div>

        <!-- 分页 -->
        <div style="margin-top: 16px; display: flex; justify-content: center;">
          <el-pagination
            v-model:current-page="pagination.page"
            v-model:page-size="pagination.size"
            :total="logs.length"
            :page-sizes="[50, 100, 200, 500]"
            layout="total, sizes, prev, pager, next, jumper"
            @current-change="handlePageChange"
            @size-change="handleSizeChange"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Refresh, Delete, User } from '@element-plus/icons-vue'
import axios from 'axios'

// 数据
const loading = ref(false)
const logs = ref([])

// 筛选器
const filters = reactive({
  method: '',
  status: '',
  username: '',
  clientIp: ''
})

// 分页
const pagination = reactive({
  page: 1,
  size: 100
})

// 计算属性 - 筛选后的日志
const filteredLogs = computed(() => {
  let result = logs.value

  if (filters.method) {
    result = result.filter(log => log.method === filters.method)
  }
  
  if (filters.status) {
    result = result.filter(log => log.status_code.toString() === filters.status)
  }
  
  if (filters.username) {
    result = result.filter(log => 
      log.username && log.username.toLowerCase().includes(filters.username.toLowerCase())
    )
  }
  
  if (filters.clientIp) {
    result = result.filter(log => 
      log.client_ip && log.client_ip.includes(filters.clientIp)
    )
  }

  // 分页
  const start = (pagination.page - 1) * pagination.size
  const end = start + pagination.size
  return result.slice(start, end)
})

// 方法
const loadData = async () => {
  loading.value = true
  try {
    const response = await axios.get('/api/logs?limit=1000')
    logs.value = response.data || []
    
    // 按时间降序排序
    logs.value.sort((a, b) => new Date(b.created_at) - new Date(a.created_at))
  } catch (error) {
    ElMessage.error('加载访问日志失败')
    console.error(error)
  } finally {
    loading.value = false
  }
}

const handleFilterChange = () => {
  pagination.page = 1  // 重置到第一页
}

const handlePageChange = (page) => {
  pagination.page = page
}

const handleSizeChange = (size) => {
  pagination.size = size
  pagination.page = 1
}

const handleClearLogs = async () => {
  try {
    await ElMessageBox.confirm(
      '确认清空所有访问日志吗？此操作不可恢复。',
      '确认清空',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await axios.delete('/api/logs')
    ElMessage.success('日志清空成功')
    loadData()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('清空失败')
      console.error(error)
    }
  }
}

const getMethodTagType = (method) => {
  const types = {
    'GET': '',
    'POST': 'success',
    'PUT': 'warning',
    'DELETE': 'danger',
    'HEAD': 'info'
  }
  return types[method] || ''
}

const getStatusTagType = (status) => {
  if (status >= 200 && status < 300) return 'success'
  if (status >= 300 && status < 400) return 'warning'
  if (status >= 400 && status < 500) return 'danger'
  if (status >= 500) return 'danger'
  return ''
}

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleString('zh-CN')
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.filter-container {
  margin-bottom: 16px;
  padding: 16px;
  background: #f8f9fa;
  border-radius: 8px;
}

.el-table .el-tag {
  font-family: monospace;
}
</style>
