import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores'

const router = createRouter({
    // 配置模式：webHistory or webHashHistory
    // import.meta.env.BASE_URL代表路径前缀的环境变量，默认是'/'
    history: createWebHistory(import.meta.env.BASE_URL),
    // 路由规则
    routes: [
        {
            path: '/',
            component: () => import('@/pages/layout/LayoutContainer.vue'), // 一级页面
            redirect: '/home',
            children: [
                { path: '/home', component: () => import('@/pages/home/HomePage.vue') },
                { path: '/problems', component: () => import('@/pages/problem/ProblemsPage.vue') },
                { path: '/problem/:id', component: () => import('@/pages/problem/ProblemPage.vue') },
                { path: '/training', component: () => import('@/pages/training/TrainingPage.vue') },
                { path: '/contest', component: () => import('@/pages/contest/ContestPage.vue') },
                { path: '/judge', component: () => import('@/pages/judge/JudgePage.vue') },
                { path: '/rank', component: () => import('@/pages/rank/RankPage.vue') },
                { path: '/discussion', component: () => import('@/pages/discussion/DiscussionPage.vue') },
                { path: '/user/profile', component: () => import('@/pages/user/UserProfilePage.vue') },
                { path: '/user/password', component: () => import('@/pages/user/UserPasswordPage.vue') },
                {
                    path: '/admin',
                    component: () => import('@/pages/cms/CmsContainer.vue'), // 后台管理（嵌套在主布局中）
                    redirect: '/admin/user',
                    children: [
                        { path: '/admin/user', component: () => import('@/pages/cms/UserManagePage.vue') },
                        { path: '/admin/problem', component: () => import('@/pages/cms/ProblemManagePage.vue') },
                        { path: '/admin/contest', component: () => import('@/pages/cms/ContestManagePage.vue') },
                        { path: '/admin/sandbox', component: () => import('@/pages/cms/SandboxManagePage.vue') },
                        { path: '/admin/notice', component: () => import('@/pages/cms/NoticeManagePage.vue') },
                    ]
                }
            ]
        },
    ]
})

/**
 * 全局权限校验，每次切换页面时都会执行
 * 访问系统直接进入首页，仅管理员路由需要登录态校验
 */
router.beforeEach((to) => {
    const userStore = useUserStore()
    if (to.path.startsWith('/admin') && userStore.userInfo.role !== 1) {
        ElMessage.error('权限不足')
        return '/'
    }
})


export default router
