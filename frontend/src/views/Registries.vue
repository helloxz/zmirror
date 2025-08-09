<template>
  <div>
    <!-- 页面头部 -->
    <div class="page-header">
      <h2 class="page-title">镜像源管理</h2>
      <p class="page-description">管理上游Docker镜像源配置</p>
    </div>

    <!-- 内容卡片 -->
    <div class="content-card">
      <div class="card-header">
        <h3 class="card-title">镜像源列表</h3>
        <el-button type="primary" @click="openDialog()">
          <el-icon><Plus /></el-icon>
          添加镜像源
        </el-button>
      </div>
      
      <div class="card-body">
        <!-- 表格 -->
        <div class="table-container">
          <el-table
            v-loading="loading"
            :data="registries"
            stripe
            style="width: 100%"
            empty-text="暂无镜像源数据"
          >
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="url" label="镜像源URL" min-width="300" />
            <el-table-column prop="priority" label="优先级" width="100" align="center" />
            <el-table-column prop="enabled" label="状态" width="100" align="center">
              <template #default="{ row }">
                <el-tag :type="row.enabled ? 'success' : 'danger'" size="small">
                  {{ row.enabled ? '启用' : '禁用' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="180" align="center">
              <template #default="{ row }">
                <el-button size="small" @click="openDialog(row)">
                  <el-icon><Edit /></el-icon>
                  编辑
                </el-button>
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

    <!-- 添加/编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
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
        <el-form-item label="URL" prop="url">
          <el-input
            v-model="form.url"
            placeholder="请输入镜像源URL，例如：https://registry-1.docker.io"
          />
        </el-form-item>
        <el-form-item label="优先级" prop="priority">
          <el-input-number
            v-model="form.priority"
            :min="1"
            :max="100"
            style="width: 100%"
          />
          <div style="font-size: 12px; color: #909399; margin-top: 4px;">
            数值越小优先级越高
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
import { Plus, Edit, Delete } from '@element-plus/icons-vue'
import axios from 'axios'

// 数据
const loading = ref(false)
const registries = ref([])

// 对话框
const dialogVisible = ref(false)
const dialogTitle = ref('')
const submitting = ref(false)

// 表单
const formRef = ref()
const form = reactive({
  id: null,
  url: '',
  priority: 1,
  enabled: true
})

const rules = {
  url: [
    { required: true, message: '请输入镜像源URL', trigger: 'blur' },
    { type: 'url', message: '请输入有效的URL', trigger: 'blur' }
  ],
  priority: [
    { required: true, message: '请输入优先级', trigger: 'blur' },
    { type: 'number', min: 1, max: 100, message: '优先级必须在1-100之间', trigger: 'blur' }
  ]
}

// 方法
const loadData = async () => {
  loading.value = true
  try {
    const response = await axios.get('/api/registries')
    registries.value = response.data || []
  } catch (error) {
    ElMessage.error('加载镜像源失败')
    console.error(error)
  } finally {
    loading.value = false
  }
}

const openDialog = (row = null) => {
  if (row) {
    // 编辑
    dialogTitle.value = '编辑镜像源'
    Object.assign(form, row)
  } else {
    // 新增
    dialogTitle.value = '添加镜像源'
    Object.assign(form, {
      id: null,
      url: '',
      priority: 1,
      enabled: true
    })
  }
  dialogVisible.value = true
}

const handleSubmit = async () => {
  if (!formRef.value) return
  
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  submitting.value = true
  
  try {
    const isEdit = form.id !== null
    const method = isEdit ? 'put' : 'post'
    
    await axios[method]('/api/registries', form)
    
    ElMessage.success(isEdit ? '更新成功' : '添加成功')
    dialogVisible.value = false
    loadData()
  } catch (error) {
    ElMessage.error('保存失败')
    console.error(error)
  } finally {
    submitting.value = false
  }
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确认删除镜像源 "${row.url}" 吗？`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await axios.delete(`/api/registries/${row.id}`)
    ElMessage.success('删除成功')
    loadData()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
      console.error(error)
    }
  }
}

onMounted(() => {
  loadData()
})
</script>
