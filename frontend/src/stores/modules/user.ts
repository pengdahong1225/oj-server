import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { User } from '@/types/user'

/**
 * 用户模块
 * token
 */
export const useUserStore = defineStore('user', () => {
    const token = ref('')
    const setToken = (value: string) => {
        token.value = value
    }
    const clearToken = () => {
        token.value = ''
    }

    const userInfo = ref<User>()
    const setUserInfo = (value: User) => {
        userInfo.value = value
    }
    const clearUserInfo = () => {
        userInfo.value = undefined
    }

    return {
        token,
        setToken,
        clearToken,
        userInfo,
        setUserInfo,
        clearUserInfo
    }
}, {
    persist: true
})