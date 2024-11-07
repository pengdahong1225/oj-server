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
      component: () => import('@/views/layout/LayoutContainer.vue'), // 一级页面
      redirect: '/home',
      children: [
        { path: '/home', component: () => import('@/views/home/HomePage.vue') },
        { path: '/problems', component: () => import('@/views/problem/ProblemsPage.vue') },
        { path: '/problem/:id', component: () => import('@/views/problem/ProblemPage.vue') },
        { path: '/contest', component: () => import('@/views/contest/ContestPage.vue') },
        { path: '/status', component: () => import('@/views/status/StatusListPage.vue') },
        { path: '/status/:id', component: () => import('@/views/status/StatusPage.vue') },
        { path: '/user/profile', component: () => import('@/views/user/UserProfilePage.vue') },
        { path: '/user/password', component: () => import('@/views/user/UserPasswordPage.vue') }
      ]
    },
    {
      path: '/manage',
      component: () => import('@/views/cms/CmsContainer.vue'), // 一级页面
      redirect: '/manage/user',
      children: [
        { path: '/manage/user', component: () => import('@/views/cms/UserManagePage.vue') },
        { path: '/manage/problem', component: () => import('@/views/cms/ProblemManagePage.vue') },
        { path: '/manage/contest', component: () => import('@/views/cms/ContestManagePage.vue') },
        { path: '/manage/template', component: () => import('@/views/cms/TemplateManagePage.vue') },
        { path: '/manage/notice', component: () => import('@/views/cms/NoticeManagePage.vue') },
      ]
    },
  ]
})

// 登录访问拦截
router.beforeEach((to) => {
  const userStore = useUserStore()
  if (!userStore.token && to.path !== '/login') {
    return '/login'
  }
  if (to.path === '/manage' && userStore.userInfo.role !== 1) {
    ElMessage.error('权限不足')
    return '/'
  }
})


export default router
