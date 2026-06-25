<template>
  <div>
    <!-- Create Restore Form -->
    <el-card style="margin-bottom: 20px">
      <template #header>
        <div style="display: flex; align-items: center">
          <el-icon style="margin-right: 8px"><Plus /></el-icon>
          <span>创建恢复</span>
        </div>
      </template>

      <el-form :model="restoreForm" label-width="120px">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="恢复名称" required>
              <el-input v-model="restoreForm.name" placeholder="例如: restore-default-20260623" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="备份名称" required>
              <el-select
                v-model="restoreForm.backupName"
                filterable
                placeholder="选择要恢复的备份"
                style="width: 100%"
                @change="handleBackupChange"
              >
                <el-option
                  v-for="backup in availableBackups"
                  :key="backup.name"
                  :label="backup.name"
                  :value="backup.name"
                >
                  <div style="display: flex; justify-content: space-between">
                    <span>{{ backup.name }}</span>
                    <el-tag v-if="backup.phase === 'Completed'" type="success" size="small">已完成</el-tag>
                  </div>
                </el-option>
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="恢复到命名空间">
              <el-input
                v-model="namespaceInput"
                placeholder="输入命名空间，多个用逗号分隔"
              />
              <div style="color: #909399; font-size: 12px; margin-top: 5px">
                留空将恢复到备份时的原命名空间
              </div>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="恢复的资源类型">
              <el-select
                v-model="restoreForm.includedResources"
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
            <el-form-item label="恢复持久卷">
              <el-switch v-model="restoreForm.restorePVs" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item>
          <el-button type="primary" @click="handleCreateRestore" :loading="creating">
            <el-icon style="margin-right: 5px"><Check /></el-icon>
            创建恢复
          </el-button>
          <el-button @click="resetForm">
            <el-icon style="margin-right: 5px"><Refresh /></el-icon>
            重置
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- Restore List -->
    <el-card>
      <template #header>
        <div style="display: flex; justify-content: space-between; align-items: center">
          <div style="display: flex; align-items: center">
            <el-icon style="margin-right: 8px"><List /></el-icon>
            <span>恢复列表</span>
          </div>
          <el-button size="small" @click="loadRestores" :loading="loading">
            <el-icon style="margin-right: 5px"><Refresh /></el-icon>
            刷新
          </el-button>
        </div>
      </template>

      <el-table :data="restores" v-loading="loading" stripe>
        <el-table-column prop="name" label="恢复名称" width="250" />
        <el-table-column prop="namespace" label="命名空间" width="120" />
        <el-table-column prop="backupName" label="来源备份" width="250" />
        <el-table-column label="状态" width="120">
          <template #default="{ row }">
            <el-tag v-if="row.phase === 'Completed'" type="success">已完成</el-tag>
            <el-tag v-else-if="row.phase === 'InProgress'" type="warning">进行中</el-tag>
            <el-tag v-else-if="row.phase === 'Failed'" type="danger">失败</el-tag>
            <el-tag v-else-if="row.phase === 'PartiallyFailed'" type="warning">部分失败</el-tag>
            <el-tag v-else type="info">{{ row.phase || '新建' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="错误/警告" width="120">
          <template #default="{ row }">
            <span style="color: #f56c6c">{{ row.errors }}</span> /
            <span style="color: #e6a23c">{{ row.warnings }}</span>
          </template>
        </el-table-column>
        <el-table-column label="开始时间" width="180">
          <template #default="{ row }">
            {{ formatTime(row.startTimestamp) }}
          </template>
        </el-table-column>
        <el-table-column label="完成时间" width="180">
          <template #default="{ row }">
            {{ formatTime(row.completionTimestamp) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="100" fixed="right">
          <template #default="{ row }">
            <el-button size="small" type="primary" link @click="viewDetails(row)">
              详情
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- Details Dialog -->
    <el-dialog v-model="detailsVisible" title="恢复详情" width="800px">
      <el-descriptions :column="2" border v-if="selectedRestore">
        <el-descriptions-item label="名称">{{ selectedRestore.name }}</el-descriptions-item>
        <el-descriptions-item label="命名空间">{{ selectedRestore.namespace }}</el-descriptions-item>
        <el-descriptions-item label="备份名称" :span="2">{{ selectedRestore.backupName }}</el-descriptions-item>
        <el-descriptions-item label="状态" :span="2">
          <el-tag v-if="selectedRestore.phase === 'Completed'" type="success">已完成</el-tag>
          <el-tag v-else-if="selectedRestore.phase === 'InProgress'" type="warning">进行中</el-tag>
          <el-tag v-else-if="selectedRestore.phase === 'Failed'" type="danger">失败</el-tag>
          <el-tag v-else-if="selectedRestore.phase === 'PartiallyFailed'" type="warning">部分失败</el-tag>
          <el-tag v-else type="info">{{ selectedRestore.phase || '新建' }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="开始时间">{{ formatTime(selectedRestore.startTimestamp) }}</el-descriptions-item>
        <el-descriptions-item label="完成时间">{{ formatTime(selectedRestore.completionTimestamp) }}</el-descriptions-item>
        <el-descriptions-item label="错误数">{{ selectedRestore.errors }}</el-descriptions-item>
        <el-descriptions-item label="警告数">{{ selectedRestore.warnings }}</el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import api from '../api/velero'

const restoreForm = ref({
  name: '',
  backupName: '',
  includedNamespaces: [],
  excludedNamespaces: [],
  includedResources: [],
  excludedResources: [],
  restorePVs: false
})

const namespaceInput = ref('')

const restores = ref([])
const availableBackups = ref([])
const loading = ref(false)
const creating = ref(false)
const detailsVisible = ref(false)
const selectedRestore = ref(null)

const loadRestores = async () => {
  loading.value = true
  try {
    const res = await api.listRestores()
    restores.value = res.data || []
  } catch (error) {
    ElMessage.error('加载恢复列表失败: ' + error.message)
  } finally {
    loading.value = false
  }
}

const loadBackups = async () => {
  try {
    const res = await api.listBackups()
    availableBackups.value = res.data || []
  } catch (error) {
    ElMessage.error('加载备份列表失败: ' + error.message)
  }
}

const handleBackupChange = (backupName) => {
  // Auto-generate restore name
  if (!restoreForm.value.name) {
    const timestamp = new Date().toISOString().slice(0, 19).replace(/[:-]/g, '').replace('T', '-')
    restoreForm.value.name = `restore-${backupName}-${timestamp}`
  }
}

const handleCreateRestore = async () => {
  if (!restoreForm.value.name) {
    ElMessage.warning('请输入恢复名称')
    return
  }
  if (!restoreForm.value.backupName) {
    ElMessage.warning('请选择要恢复的备份')
    return
  }

  creating.value = true
  try {
    const payload = { ...restoreForm.value }

    // Parse namespace input (comma-separated)
    if (namespaceInput.value.trim()) {
      payload.includedNamespaces = namespaceInput.value.split(',').map(s => s.trim()).filter(s => s)
    } else {
      payload.includedNamespaces = []
    }

    // Remove empty fields
    if (payload.includedNamespaces.length === 0) delete payload.includedNamespaces
    if (payload.excludedNamespaces.length === 0) delete payload.excludedNamespaces
    if (payload.includedResources.length === 0) delete payload.includedResources
    if (payload.excludedResources.length === 0) delete payload.excludedResources

    await api.createRestore(payload)
    ElMessage.success('恢复创建成功')
    resetForm()
    loadRestores()
  } catch (error) {
    ElMessage.error('创建恢复失败: ' + error.message)
  } finally {
    creating.value = false
  }
}

const viewDetails = (restore) => {
  selectedRestore.value = restore
  detailsVisible.value = true
}

const resetForm = () => {
  restoreForm.value = {
    name: '',
    backupName: '',
    includedNamespaces: [],
    excludedNamespaces: [],
    includedResources: [],
    excludedResources: [],
    restorePVs: false
  }
  namespaceInput.value = ''
}

const formatTime = (time) => {
  if (!time) return '-'
  return new Date(time).toLocaleString('zh-CN')
}

onMounted(() => {
  loadRestores()
  loadBackups()
})
</script>
