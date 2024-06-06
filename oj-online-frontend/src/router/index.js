import Vue from 'vue'
import VueRouter from 'vue-router'

/**
 * 默认加载
 */
import Layout from '@/views/layout'
import Home from '@/views/layout/home'
import ProblemSet from '@/views/layout/prolist'
import Rank from '@/views/layout/rank'
import User from '@/views/layout/user'
/**
 * 懒加载 => 异步组件
 */
const Login = () => import('@/views/login')
const Register = () => import('@/views/register')
const Detail = () => import('@/views/prodetail')

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
        path: '/problemset', component: ProblemSet
      },
      {
        path: '/rank', component: Rank
      },
      {
        path: '/user', component: User
      }
    ]
  },
  { path: '/login', component: Login },
  { path: '/register', component: Register },
  { path: '/detail', component: Detail }
]

const router = new VueRouter({
  routes
})

// const authUrls = ['/user']

// /**
//  * 全局前置导航守卫
//  * to: 即将进入的目标路由对象
//  * from: 当前导航正要离开的路由
//  * next: 是否放行
//  * next(): 放行，继续导航
//  * next(路径): 拦截当前导航，拦截到next里面配置的路径
//  */
// router.beforeEach((to, from, next) => {
//   if (!authUrls.includes(to.path)) {
//     next()
//   } else {
//     const token = store.getters.token
//     console.log(token)
//     if (token) {
//       next()
//     } else {
//       next('/login')
//     }
//   }
// })

export default router
