<template>
  <div>
    <!-- Create Schedule Form -->
    <el-card style="margin-bottom: 20px">
      <template #header>
        <div style="display: flex; align-items: center">
          <el-icon style="margin-right: 8px"><Plus /></el-icon>
          <span>创建定时备份任务</span>
        </div>
      </template>

      <el-form :model="scheduleForm" label-width="140px">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="任务名称" required>
              <el-input v-model="scheduleForm.name" placeholder="例如: daily-backup" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="Cron 表达式" required>
              <el-select
                v-model="scheduleForm.schedule"
                filterable
                allow-create
                placeholder="选择或自定义"
                style="width: 100%"
              >
                <el-option label="每小时 (0 * * * *)" value="0 * * * *" />
                <el-option label="每天凌晨2点 (0 2 * * *)" value="0 2 * * *" />
                <el-option label="每天中午12点 (0 12 * * *)" value="0 12 * * *" />
                <el-option label="每周日凌晨2点 (0 2 * * 0)" value="0 2 * * 0" />
                <el-option label="每月1日凌晨2点 (0 2 1 * *)" value="0 2 1 * *" />
              </el-select>
              <div style="color: #909399; font-size: 12px; margin-top: 5px">
                格式: 分 时 日 月 周 (例如: 0 2 * * * 表示每天凌晨2点)
              </div>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="命名空间">
              <el-input
                v-model="namespaceInput"
                placeholder="输入命名空间，多个用逗号分隔"
              />
              <div style="color: #909399; font-size: 12px; margin-top: 5px">
                留空表示备份所有命名空间
              </div>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="资源类型">
              <el-select
                v-model="scheduleForm.includedResources"
                multiple
                filterable
                placeholder="选择资源类型，不选表示全部"
                style="width: 100%"
              >
                <el-option label="Deployment" value="deployments.apps" />
                <el-option label="Service" value="services" />
                <el-option label="ConfigMap" value="configmaps" />
                <el-option label="Secret" value="secrets" />
                <el-option label="PersistentVolumeClaim" value="persistentvolumeclaims" />
                <el-option label="StatefulSet" value="statefulsets.apps" />
                <el-option label="DaemonSet" value="daemonsets.apps" />
                <el-option label="Pod" value="pods" />
                <el-option label="Ingress" value="ingresses.networking.k8s.io" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="保留时间 (TTL)">
              <el-input v-model="scheduleForm.ttl" placeholder="例如: 720h (30天)">
                <template #append>小时</template>
              </el-input>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="存储位置">
              <el-select
                v-model="scheduleForm.storageLocation"
                filterable
                placeholder="选择存储位置（默认使用第一个）"
                style="width: 100%"
              >
                <el-option
                  v-for="location in storageLocations"
                  :key="location.name"
                  :label="location.name"
                  :value="location.name"
                >
                  <div style="display: flex; justify-content: space-between">
                    <span>{{ location.name }}</span>
                    <el-tag v-if="location.isDefault" type="success" size="small">默认</el-tag>
                  </div>
                </el-option>
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item>
          <el-button type="primary" @click="handleCreateSchedule" :loading="creating">
            <el-icon style="margin-right: 5px"><Check /></el-icon>
            创建任务
          </el-button>
          <el-button @click="resetForm">
            <el-icon style="margin-right: 5px"><Refresh /></el-icon>
            重置
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- Schedule List -->
    <el-card>
      <template #header>
        <div style="display: flex; justify-content: space-between; align-items: center">
          <div style="display: flex; align-items: center">
            <el-icon style="margin-right: 8px"><List /></el-icon>
            <span>定时任务列表</span>
          </div>
          <el-button size="small" @click="loadSchedules" :loading="loading">
            <el-icon style="margin-right: 5px"><Refresh /></el-icon>
            刷新
          </el-button>
        </div>
      </template>

      <el-table :data="schedules" v-loading="loading" stripe>
        <el-table-column prop="name" label="任务名称" width="200" />
        <el-table-column prop="namespace" label="命名空间" width="120" />
        <el-table-column prop="schedule" label="Cron 表达式" width="150" />
        <el-table-column label="状态" width="120">
          <template #default="{ row }">
            <el-tag v-if="row.phase === 'Enabled'" type="success">启用</el-tag>
            <el-tag v-else-if="row.phase === 'FailedValidation'" type="danger">验证失败</el-tag>
            <el-tag v-else type="info">{{ row.phase || '未知' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="storageLocation" label="存储位置" width="120" />
        <el-table-column prop="ttl" label="保留时间" width="120" />
        <el-table-column label="最后备份时间" width="180">
          <template #default="{ row }">
            {{ formatTime(row.lastBackup) }}
          </template>
        </el-table-column>
        <el-table-column label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatTime(row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button size="small" type="primary" link @click="viewDetails(row)">
              详情
            </el-button>
            <el-popconfirm
              title="确定要删除这个定时任务吗？"
              @confirm="handleDeleteSchedule(row.name)"
            >
              <template #reference>
                <el-button size="small" type="danger" link>删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- Details Dialog -->
    <el-dialog v-model="detailsVisible" title="定时任务详情" width="800px">
      <el-descriptions :column="2" border v-if="selectedSchedule">
        <el-descriptions-item label="名称">{{ selectedSchedule.name }}</el-descriptions-item>
        <el-descriptions-item label="命名空间">{{ selectedSchedule.namespace }}</el-descriptions-item>
        <el-descriptions-item label="Cron 表达式">{{ selectedSchedule.schedule }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag v-if="selectedSchedule.phase === 'Enabled'" type="success">启用</el-tag>
          <el-tag v-else-if="selectedSchedule.phase === 'FailedValidation'" type="danger">验证失败</el-tag>
          <el-tag v-else type="info">{{ selectedSchedule.phase || '未知' }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="存储位置">{{ selectedSchedule.storageLocation }}</el-descriptions-item>
        <el-descriptions-item label="保留时间">{{ selectedSchedule.ttl }}</el-descriptions-item>
        <el-descriptions-item label="最后备份时间" :span="2">
          {{ formatTime(selectedSchedule.lastBackup) }}
        </el-descriptions-item>
        <el-descriptions-item label="创建时间" :span="2">
          {{ formatTime(selectedSchedule.createdAt) }}
        </el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import api from '../api/velero'

const scheduleForm = ref({
  name: '',
  schedule: '',
  includedNamespaces: [],
  excludedNamespaces: [],
  includedResources: [],
  excludedResources: [],
  storageLocation: '',
  ttl: '720h'
})

const namespaceInput = ref('')
const storageLocations = ref([])

const schedules = ref([])
const loading = ref(false)
const creating = ref(false)
const detailsVisible = ref(false)
const selectedSchedule = ref(null)

const loadSchedules = async () => {
  loading.value = true
  try {
    const res = await api.listSchedules()
    schedules.value = res.data || []
  } catch (error) {
    ElMessage.error('加载定时任务列表失败: ' + error.message)
  } finally {
    loading.value = false
  }
}

const loadStorageLocations = async () => {
  try {
    const res = await api.listStorageLocations()
    storageLocations.value = res.data || []
  } catch (error) {
    ElMessage.error('加载存储位置列表失败: ' + error.message)
  }
}

const handleCreateSchedule = async () => {
  if (!scheduleForm.value.name) {
    ElMessage.warning('请输入任务名称')
    return
  }
  if (!scheduleForm.value.schedule) {
    ElMessage.warning('请输入 Cron 表达式')
    return
  }

  creating.value = true
  try {
    const payload = { ...scheduleForm.value }

    // Parse namespace input (comma-separated)
    if (namespaceInput.value.trim()) {
      payload.includedNamespaces = namespaceInput.value.split(',').map(s => s.trim()).filter(s => s)
    } else {
      payload.includedNamespaces = []
    }

    // Remove empty fields
    if (!payload.storageLocation) delete payload.storageLocation
    if (!payload.ttl) delete payload.ttl
    if (payload.excludedNamespaces.length === 0) delete payload.excludedNamespaces
    if (payload.excludedResources.length === 0) delete payload.excludedResources
    if (payload.includedResources.length === 0) delete payload.includedResources

    await api.createSchedule(payload)
    ElMessage.success('定时任务创建成功')
    resetForm()
    loadSchedules()
  } catch (error) {
    ElMessage.error('创建定时任务失败: ' + error.message)
  } finally {
    creating.value = false
  }
}

const handleDeleteSchedule = async (name) => {
  try {
    await api.deleteSchedule(name)
    ElMessage.success('删除成功')
    loadSchedules()
  } catch (error) {
    ElMessage.error('删除失败: ' + error.message)
  }
}

const viewDetails = (schedule) => {
  selectedSchedule.value = schedule
  detailsVisible.value = true
}

const resetForm = () => {
  scheduleForm.value = {
    name: '',
    schedule: '',
    includedNamespaces: [],
    excludedNamespaces: [],
    includedResources: [],
    excludedResources: [],
    storageLocation: '',
    ttl: '720h'
  }
  namespaceInput.value = ''
}

const formatTime = (time) => {
  if (!time) return '-'
  return new Date(time).toLocaleString('zh-CN')
}

onMounted(() => {
  loadSchedules()
  loadStorageLocations()
})
</script>
