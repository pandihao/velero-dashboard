<template>
  <div>
    <!-- Create Backup Form -->
    <el-card style="margin-bottom: 20px">
      <template #header>
        <div style="display: flex; align-items: center">
          <el-icon style="margin-right: 8px"><Plus /></el-icon>
          <span>创建备份</span>
        </div>
      </template>

      <el-form :model="backupForm" label-width="120px">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="备份名称" required>
              <el-input v-model="backupForm.name" placeholder="例如: backup-default-20260623" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="存储位置">
              <el-select
                v-model="backupForm.storageLocation"
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

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="命名空间">
              <el-input
                v-model="namespaceInput"
                placeholder="输入命名空间，多个用逗号分隔，如: default,kube-system"
              />
              <div style="color: #909399; font-size: 12px; margin-top: 5px">
                留空表示备份所有命名空间
              </div>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="资源类型">
              <el-select
                v-model="backupForm.includedResources"
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
              <el-input v-model="backupForm.ttl" placeholder="例如: 720h (30天)">
                <template #append>小时</template>
              </el-input>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="快照卷">
              <el-switch v-model="backupForm.snapshotVolumes" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item>
          <el-button type="primary" @click="handleCreateBackup" :loading="creating">
            <el-icon style="margin-right: 5px"><Check /></el-icon>
            创建备份
          </el-button>
          <el-button @click="resetForm">
            <el-icon style="margin-right: 5px"><Refresh /></el-icon>
            重置
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- Backup List -->
    <el-card>
      <template #header>
        <div style="display: flex; justify-content: space-between; align-items: center">
          <div style="display: flex; align-items: center">
            <el-icon style="margin-right: 8px"><List /></el-icon>
            <span>备份列表</span>
          </div>
          <el-button size="small" @click="loadBackups" :loading="loading">
            <el-icon style="margin-right: 5px"><Refresh /></el-icon>
            刷新
          </el-button>
        </div>
      </template>

      <el-table :data="backups" v-loading="loading" stripe>
        <el-table-column prop="name" label="备份名称" width="250" />
        <el-table-column prop="namespace" label="命名空间" width="120" />
        <el-table-column label="状态" width="120">
          <template #default="{ row }">
            <el-tag v-if="row.phase === 'Completed'" type="success">已完成</el-tag>
            <el-tag v-else-if="row.phase === 'InProgress'" type="warning">进行中</el-tag>
            <el-tag v-else-if="row.phase === 'Failed'" type="danger">失败</el-tag>
            <el-tag v-else type="info">{{ row.phase || '新建' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="包含的命名空间" width="200">
          <template #default="{ row }">
            <el-tag v-for="ns in row.includedNamespaces" :key="ns" size="small" style="margin-right: 5px">
              {{ ns }}
            </el-tag>
            <span v-if="!row.includedNamespaces || row.includedNamespaces.length === 0">全部</span>
          </template>
        </el-table-column>
        <el-table-column prop="storageLocation" label="存储位置" width="120" />
        <el-table-column prop="ttl" label="保留时间" width="120" />
        <el-table-column label="错误/警告" width="120">
          <template #default="{ row }">
            <span style="color: #f56c6c">{{ row.errors }}</span> /
            <span style="color: #e6a23c">{{ row.warnings }}</span>
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
              title="确定要删除这个备份吗？"
              @confirm="handleDeleteBackup(row.name)"
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
    <el-dialog v-model="detailsVisible" title="备份详情" width="800px">
      <el-descriptions :column="2" border v-if="selectedBackup">
        <el-descriptions-item label="名称">{{ selectedBackup.name }}</el-descriptions-item>
        <el-descriptions-item label="命名空间">{{ selectedBackup.namespace }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag v-if="selectedBackup.phase === 'Completed'" type="success">已完成</el-tag>
          <el-tag v-else-if="selectedBackup.phase === 'InProgress'" type="warning">进行中</el-tag>
          <el-tag v-else-if="selectedBackup.phase === 'Failed'" type="danger">失败</el-tag>
          <el-tag v-else type="info">{{ selectedBackup.phase || '新建' }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="存储位置">{{ selectedBackup.storageLocation }}</el-descriptions-item>
        <el-descriptions-item label="包含的命名空间" :span="2">
          {{ selectedBackup.includedNamespaces?.join(', ') || '全部' }}
        </el-descriptions-item>
        <el-descriptions-item label="排除的命名空间" :span="2">
          {{ selectedBackup.excludedNamespaces?.join(', ') || '无' }}
        </el-descriptions-item>
        <el-descriptions-item label="保留时间">{{ selectedBackup.ttl }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ formatTime(selectedBackup.createdAt) }}</el-descriptions-item>
        <el-descriptions-item label="开始时间">{{ formatTime(selectedBackup.startTimestamp) }}</el-descriptions-item>
        <el-descriptions-item label="完成时间">{{ formatTime(selectedBackup.completionTimestamp) }}</el-descriptions-item>
        <el-descriptions-item label="错误数">{{ selectedBackup.errors }}</el-descriptions-item>
        <el-descriptions-item label="警告数">{{ selectedBackup.warnings }}</el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import api from '../api/velero'

const backupForm = ref({
  name: '',
  includedNamespaces: [],
  excludedNamespaces: [],
  includedResources: [],
  excludedResources: [],
  storageLocation: '',
  ttl: '720h',
  snapshotVolumes: false
})

const namespaceInput = ref('')
const storageLocations = ref([])

const backups = ref([])
const loading = ref(false)
const creating = ref(false)
const detailsVisible = ref(false)
const selectedBackup = ref(null)

const loadBackups = async () => {
  loading.value = true
  try {
    const res = await api.listBackups()
    backups.value = res.data || []
  } catch (error) {
    ElMessage.error('加载备份列表失败: ' + error.message)
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

const handleCreateBackup = async () => {
  if (!backupForm.value.name) {
    ElMessage.warning('请输入备份名称')
    return
  }

  creating.value = true
  try {
    const payload = { ...backupForm.value }

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

    await api.createBackup(payload)
    ElMessage.success('备份创建成功')
    resetForm()
    loadBackups()
  } catch (error) {
    ElMessage.error('创建备份失败: ' + error.message)
  } finally {
    creating.value = false
  }
}

const handleDeleteBackup = async (name) => {
  try {
    await api.deleteBackup(name)
    ElMessage.success('删除成功')
    loadBackups()
  } catch (error) {
    ElMessage.error('删除失败: ' + error.message)
  }
}

const viewDetails = (backup) => {
  selectedBackup.value = backup
  detailsVisible.value = true
}

const resetForm = () => {
  backupForm.value = {
    name: '',
    includedNamespaces: [],
    excludedNamespaces: [],
    includedResources: [],
    excludedResources: [],
    storageLocation: '',
    ttl: '720h',
    snapshotVolumes: false
  }
  namespaceInput.value = ''
}

const formatTime = (time) => {
  if (!time) return '-'
  return new Date(time).toLocaleString('zh-CN')
}

onMounted(() => {
  loadBackups()
  loadStorageLocations()
})
</script>
