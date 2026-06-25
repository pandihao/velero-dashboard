<template>
  <div id="app">
    <el-container style="height: 100vh">
      <el-aside width="200px" style="background-color: #304156">
        <div style="padding: 20px; text-align: center; color: #fff; font-size: 18px; font-weight: bold">
          Velero Console
        </div>
        <el-menu
          :default-active="$route.path"
          router
          background-color="#304156"
          text-color="#bfcbd9"
          active-text-color="#409EFF"
        >
          <el-menu-item index="/backup">
            <el-icon><Document /></el-icon>
            <span>备份管理</span>
          </el-menu-item>
          <el-menu-item index="/restore">
            <el-icon><RefreshRight /></el-icon>
            <span>恢复管理</span>
          </el-menu-item>
          <el-menu-item index="/schedule">
            <el-icon><Clock /></el-icon>
            <span>定时备份</span>
          </el-menu-item>
        </el-menu>
      </el-aside>

      <el-container>
        <el-header style="background-color: #fff; box-shadow: 0 1px 4px rgba(0,21,41,.08); display: flex; justify-content: space-between; align-items: center">
          <div style="font-size: 20px; font-weight: 500">
            {{ pageTitle }}
          </div>
          <div style="display: flex; align-items: center; gap: 10px">
            <span style="color: #909399; font-size: 14px">集群:</span>
            <el-select
              v-model="selectedCluster"
              @change="handleClusterChange"
              style="width: 200px"
              size="default"
            >
              <el-option
                v-for="cluster in clusters"
                :key="cluster"
                :label="cluster"
                :value="cluster"
              />
            </el-select>
          </div>
        </el-header>

        <el-main style="background-color: #f0f2f5">
          <router-view :key="selectedCluster" />
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import api from './api/velero'

const route = useRoute()

const pageTitle = computed(() => {
  const titles = {
    '/backup': '备份管理',
    '/restore': '恢复管理',
    '/schedule': '定时备份任务'
  }
  return titles[route.path] || 'Velero Console'
})

const clusters = ref([])
const selectedCluster = ref('')

const loadClusters = async () => {
  try {
    const res = await api.listClusters()
    clusters.value = res.data.clusters || []
    selectedCluster.value = res.data.default || clusters.value[0] || ''
    if (selectedCluster.value) {
      api.setCluster(selectedCluster.value)
    }
  } catch (error) {
    ElMessage.error('加载集群列表失败: ' + error.message)
  }
}

const handleClusterChange = (cluster) => {
  api.setCluster(cluster)
  ElMessage.success(`已切换到集群: ${cluster}`)
}

onMounted(() => {
  loadClusters()
})
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
}

#app {
  width: 100%;
  height: 100vh;
}

.el-aside {
  overflow-y: auto;
}

.el-main {
  padding: 20px;
  overflow-y: auto;
}
</style>
