import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { User } from '@/types/user'

/**
 * 用户模块
 * token
 */
export const useUserStore = defineStore('user', () => {
    const userInfo = ref<User>({
        id: 0,
        nickname: '',
        avatar_url: '',
        mobile: 0,
        email: '',
        role: 0,
        token: '',
        gender: 0
    })
    const setUserInfo = (value: User) => {
        userInfo.value = value
    }
    const clearUserInfo = () => {
        userInfo.value.id = 0
        userInfo.value.nickname = ''
        userInfo.value.avatar_url = ''
        userInfo.value.mobile = 0
        userInfo.value.email = ''
        userInfo.value.role = 0
        userInfo.value.token = ''
        userInfo.value.gender = 0
    }

    return {
        userInfo,
        setUserInfo,
        clearUserInfo
    }
}, {
    persist: true
})