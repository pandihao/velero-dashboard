import { createRouter, createWebHistory } from 'vue-router'
import BackupView from './views/BackupView.vue'
import RestoreView from './views/RestoreView.vue'
import ScheduleView from './views/ScheduleView.vue'

const routes = [
  {
    path: '/',
    redirect: '/backup'
  },
  {
    path: '/backup',
    name: 'Backup',
    component: BackupView
  },
  {
    path: '/restore',
    name: 'Restore',
    component: RestoreView
  },
  {
    path: '/schedule',
    name: 'Schedule',
    component: ScheduleView
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
