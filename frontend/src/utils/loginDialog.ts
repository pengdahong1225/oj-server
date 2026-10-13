import { ref } from 'vue'

// 登录框的全局显示状态：供 axios 拦截器（需要登录时）与 LayoutContainer 共享
export const loginDialogVisible = ref(false)

export function showLoginDialog() {
    loginDialogVisible.value = true
}

export function hideLoginDialog() {
    loginDialogVisible.value = false
}
