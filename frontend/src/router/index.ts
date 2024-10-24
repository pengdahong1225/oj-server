import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores'

const router = createRouter({
  // 配置模式：webHistory or webHashHistory
  // import.meta.env.BASE_URL代表路径前缀的环境变量，默认是'/'
  history: createWebHistory(import.meta.env.BASE_URL),
  // 路由规则
  routes: [
    { path: '/login', component: () => import('@/views/login/LoginPage.vue') },
    {
      path: '/',
      component: () => import('@/views/layout/LayoutContainer.vue'),
      redirect: '/problems',
      children: [
        { path: '/problems', component: () => import('@/views/problem/ProblemsPage.vue') },
        { path: '/problem/:id', component: () => import('@/views/problem/ProblemPage.vue') },
        { path: '/contest', component: () => import('@/views/contest/ContestPage.vue') },
        { path: '/status', component: () => import('@/views/status/StatusPage.vue') },
        { path: '/user/profile', component: () => import('@/views/user/UserProfilePage.vue') },
        { path: '/user/avatar', component: () => import('@/views/user/UserAvatorPage.vue') },
        { path: '/user/password', component: () => import('@/views/user/UserPasswordPage.vue') },
      ]
    }
  ]
})

// 登录访问拦截
// router.beforeEach((to) => {
//   const userStore = useUserStore()
//   if (!userStore.token && to.path !== '/login') return '/login'
// })


export default router
