<template>
  <div>
    <!-- 页面头部 -->
    <div class="page-header">
      <h2 class="page-title">白名单管理</h2>
      <p class="page-description">配置允许访问的镜像前缀白名单</p>
    </div>

    <!-- 内容卡片 -->
    <div class="content-card">
      <div class="card-header">
        <h3 class="card-title">白名单列表</h3>
        <el-button type="primary" @click="openDialog()">
          <el-icon><Plus /></el-icon>
          添加白名单
        </el-button>
      </div>
      
      <div class="card-body">
        <!-- 表格 -->
        <div class="table-container">
          <el-table
            v-loading="loading"
            :data="whitelists"
            stripe
            style="width: 100%"
            empty-text="暂无白名单数据"
          >
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="prefix" label="镜像前缀" min-width="300">
              <template #default="{ row }">
                <el-tag type="info" size="small" style="margin-right: 8px;">
                  {{ row.prefix }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="enabled" label="状态" width="100" align="center">
              <template #default="{ row }">
                <el-tag :type="row.enabled ? 'success' : 'danger'" size="small">
                  {{ row.enabled ? '启用' : '禁用' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="created_at" label="创建时间" width="180" align="center">
              <template #default="{ row }">
                {{ formatDate(row.created_at) }}
              </template>
            </el-table-column>
            <el-table-column label="操作" width="120" align="center">
              <template #default="{ row }">
                <el-button size="small" type="danger" @click="handleDelete(row)">
                  <el-icon><Delete /></el-icon>
                  删除
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </div>
    </div>

    <!-- 添加对话框 -->
    <el-dialog
      v-model="dialogVisible"
      title="添加白名单"
      width="500px"
      class="form-dialog"
      :close-on-click-modal="false"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="100px"
        class="form-container"
        label-position="top"
      >
        <el-form-item label="镜像前缀" prop="prefix">
          <el-input
            v-model="form.prefix"
            placeholder="例如：helloxz、nginx、mysql"
          />
          <div style="font-size: 12px; color: #909399; margin-top: 4px;">
            支持完整镜像名或前缀匹配，如 "helloxz" 将匹配 "helloxz/app:latest"
          </div>
        </el-form-item>
        <el-form-item label="状态">
          <el-switch
            v-model="form.enabled"
            active-text="启用"
            inactive-text="禁用"
          />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">
          保存
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Delete } from '@element-plus/icons-vue'
import axios from 'axios'

// 数据
const loading = ref(false)
const whitelists = ref([])

// 对话框
const dialogVisible = ref(false)
const submitting = ref(false)

// 表单
const formRef = ref()
const form = reactive({
  prefix: '',
  enabled: true
})

const rules = {
  prefix: [
    { required: true, message: '请输入镜像前缀', trigger: 'blur' },
    { min: 1, max: 100, message: '长度在 1 到 100 个字符', trigger: 'blur' }
  ]
}

// 方法
const loadData = async () => {
  loading.value = true
  try {
    const response = await axios.get('/api/whitelists')
    whitelists.value = response.data || []
  } catch (error) {
    ElMessage.error('加载白名单失败')
    console.error(error)
  } finally {
    loading.value = false
  }
}

const openDialog = () => {
  Object.assign(form, {
    prefix: '',
    enabled: true
  })
  dialogVisible.value = true
}

const handleSubmit = async () => {
  if (!formRef.value) return
  
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  submitting.value = true
  
  try {
    await axios.post('/api/whitelists', form)
    
    ElMessage.success('添加成功')
    dialogVisible.value = false
    loadData()
  } catch (error) {
    ElMessage.error('添加失败')
    console.error(error)
  } finally {
    submitting.value = false
  }
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确认删除白名单 "${row.prefix}" 吗？`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await axios.delete(`/api/whitelists/${row.id}`)
    ElMessage.success('删除成功')
    loadData()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
      console.error(error)
    }
  }
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
.el-tag {
  font-family: monospace;
}
</style>
