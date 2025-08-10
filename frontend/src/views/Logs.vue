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
        <!-- 表格 -->
        <div class="table-container">
          <el-table
            v-loading="loading"
            :data="paginatedLogs"
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

// 分页
const pagination = reactive({
  page: 1,
  size: 100
})

// 计算属性 - 分页后的日志
const paginatedLogs = computed(() => {
  const start = (pagination.page - 1) * pagination.size
  const end = start + pagination.size
  return logs.value.slice(start, end)
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

.el-table .el-tag {
  font-family: monospace;
}
</style>
