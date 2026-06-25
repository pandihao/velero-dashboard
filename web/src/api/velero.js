import axios from 'axios'

const api = axios.create({
  baseURL: '/api/v1',
  timeout: 30000
})

// Store current cluster
let currentCluster = ''

// Request interceptor
api.interceptors.request.use(
  config => {
    // Add cluster parameter to all requests
    if (currentCluster) {
      config.params = config.params || {}
      config.params.cluster = currentCluster
    }
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

// Response interceptor
api.interceptors.response.use(
  response => {
    return response.data
  },
  error => {
    const message = error.response?.data?.message || error.message || 'Request failed'
    return Promise.reject(new Error(message))
  }
)

export default {
  // Set current cluster
  setCluster(cluster) {
    currentCluster = cluster
  },

  // Get current cluster
  getCluster() {
    return currentCluster
  },

  // Cluster APIs
  listClusters() {
    return axios.get('/api/v1/clusters').then(res => res.data)
  },

  // Backup APIs
  createBackup(data) {
    return api.post('/backups', data)
  },
  listBackups() {
    return api.get('/backups')
  },
  getBackup(name) {
    return api.get(`/backups/${name}`)
  },
  deleteBackup(name) {
    return api.delete(`/backups/${name}`)
  },

  // Restore APIs
  createRestore(data) {
    return api.post('/restores', data)
  },
  listRestores(backupName = '') {
    return api.get('/restores', { params: backupName ? { backupName } : {} })
  },
  getRestore(name) {
    return api.get(`/restores/${name}`)
  },

  // Schedule APIs
  createSchedule(data) {
    return api.post('/schedules', data)
  },
  listSchedules() {
    return api.get('/schedules')
  },
  getSchedule(name) {
    return api.get(`/schedules/${name}`)
  },
  deleteSchedule(name) {
    return api.delete(`/schedules/${name}`)
  },

  // Storage Location APIs
  listStorageLocations() {
    return api.get('/storage-locations')
  },
  getStorageLocation(name) {
    return api.get(`/storage-locations/${name}`)
  }
}
