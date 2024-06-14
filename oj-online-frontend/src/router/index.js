import Vue from 'vue'
import VueRouter from 'vue-router'

import store from '@/store'
/**
 * 默认加载
 */
import Layout from '@/views'
import Home from '@/views/layout/home'
import ProblemList from '@/views/layout/problemlist'
import Rank from '@/views/layout/rank'

/**
 * 懒加载 => 异步组件
 */
const User = () => import('@/views/layout/user')
const Problem = () => import('@/views/layout/problem')
const Submission = () => import('@/views/layout/submission')
const Setting = () => import('@/views/layout/setting')

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    component: Layout,
    redirect: '/home',
    children: [
      {
        path: '/home', component: Home
      },
      {
        path: '/problemList', component: ProblemList
      },
      {
        path: '/rank', component: Rank
      },
      {
        path: '/user', component: User
      },
      {
        path: '/problem/:id', component: Problem
      },
      {
        path: '/submission', component: Submission
      },
      {
        path: '/setting', component: Setting
      }
    ]
  },

  { path: '*', redirect: '/home' }
]

const router = new VueRouter({
  routes
})

/**
 * 全局前置导航守卫
 * to: 即将进入的目标路由对象
 * from: 当前导航正要离开的路由
 * next: 是否放行
 * next(): 放行，继续导航
 * next(路径): 拦截当前导航，拦截到next里面配置的路径
 */
const authUrls = []
router.beforeEach((to, from, next) => {
  if (!authUrls.includes(to.path)) {
    next()
  } else {
    const token = store.getters.token
    console.log(token)
    if (token) {
      next()
    } else {
      next('/login')
    }
  }
})

// 解决push报错
const originalPush = VueRouter.prototype.push
VueRouter.prototype.push = function (location) {
  return originalPush.call(this, location).catch(err => err)
}

export default router
