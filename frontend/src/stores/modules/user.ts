import { defineStore } from 'pinia'
import { ref } from 'vue'

/**
 * 用户模块
 * token
 */
export const useUserStore = defineStore('user', () => {
    const token = ref('')
    const userInfo = ref<API.UserInfo>({
        uid: 0,
        nickname: '',
        avatarUrl: '',
        mobile: '',
        email: '',
        role: 0,
        gender: 0,
        createAt: 0,
    })
    const setUserInfo = (value: API.UserInfo) => {
        userInfo.value = value
    }
    const setToken = (value: string) => {
        token.value = value
    }
    const clearToken = () => {
        token.value = ''
    }
    const clearUserInfo = () => {
        userInfo.value.uid = 0
        userInfo.value.nickname = ''
        userInfo.value.avatarUrl = ''
        userInfo.value.mobile = ''
        userInfo.value.email = ''
        userInfo.value.role = 0
        userInfo.value.gender = 0
        userInfo.value.createAt = 0
        token.value = ''
        // 清除本地持久化的登录态，避免刷新后恢复
        localStorage.removeItem('user')
    }

    return {
        userInfo,
        token,
        setUserInfo,
        clearUserInfo,
        setToken,
        clearToken
    }
}, {
    persist: true // 持久化
})
