import { defineStore } from 'pinia'
import { ref } from 'vue'

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

    return {
        token,
        setToken,
        clearToken
    }
}, {
    persist: true
})