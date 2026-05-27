import { createRouter, createWebHashHistory,createWebHistory } from 'vue-router'

const routes = [{
    path: '/',
    name: 'PixelCard',
    component: () => import('@/view/homepage.vue')
},
{
  path: '/login',
  name: 'Login',
  component: () => import('@/view/login/login.vue')
},
{
  path: '/forgetPassword',
  name: 'forgetPassword',
  component: () => import('@/view/login/refound_password.vue')
}
,
{
  path: '/adminLogin/:code',
  name: 'adminLogin',
  component: () => import('@/view/login/admin_snap_login.vue')
}
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
