import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores'

const router = createRouter({
    // 配置模式：webHistory or webHashHistory
    // import.meta.env.BASE_URL代表路径前缀的环境变量，默认是'/'
    history: createWebHistory(import.meta.env.BASE_URL),
    // 路由规则
    routes: [
        { path: '/login', component: () => import('@/pages/login/LoginPage.vue') },
        {
            path: '/',
            component: () => import('@/pages/layout/LayoutContainer.vue'), // 一级页面
            redirect: '/home',
            children: [
                { path: '/home', component: () => import('@/pages/home/HomePage.vue') },
                { path: '/problems', component: () => import('@/pages/problem/ProblemsPage.vue') },
                { path: '/problem/:id', component: () => import('@/pages/problem/ProblemPage.vue') },
                { path: '/contest', component: () => import('@/pages/contest/ContestPage.vue') },
                { path: '/status', component: () => import('@/pages/status/StatusListPage.vue') },
                { path: '/status/:id', component: () => import('@/pages/status/StatusPage.vue') },
                { path: '/user/profile', component: () => import('@/pages/user/UserProfilePage.vue') },
                { path: '/user/password', component: () => import('@/pages/user/UserPasswordPage.vue') }
            ]
        },
        {
            path: '/manage',
            component: () => import('@/pages/cms/CmsContainer.vue'), // 一级页面
            redirect: '/manage/user',
            children: [
                { path: '/manage/user', component: () => import('@/pages/cms/UserManagePage.vue') },
                { path: '/manage/problem', component: () => import('@/pages/cms/ProblemManagePage.vue') },
                { path: '/manage/contest', component: () => import('@/pages/cms/ContestManagePage.vue') },
                { path: '/manage/template', component: () => import('@/pages/cms/TemplateManagePage.vue') },
                { path: '/manage/notice', component: () => import('@/pages/cms/NoticeManagePage.vue') },
            ]
        },
    ]
})

/**
 * 登录访问拦截
 * 全局权限校验，每次切换页面时都会执行
 */
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
