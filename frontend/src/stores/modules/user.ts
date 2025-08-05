import { defineStore } from 'pinia'
import { ref } from 'vue'

/**
 * 用户模块
 * token
 */
export const useUserStore = defineStore('user', () => {
    const token = ref('')
    const userInfo = ref<API.UserProfile>({
        uid: 0,
        nick_name: '',
        avatar_url: '',
        mobile: 0,
        email: '',
        role: 0,
        gender: 0
    })
    const setUserInfo = (value: API.UserProfile) => {
        userInfo.value = value
    }
    const setToken = (value: string) => {
        token.value = value
    }
    const clearUserInfo = () => {
        userInfo.value.uid = 0
        userInfo.value.nick_name = ''
        userInfo.value.avatar_url = ''
        userInfo.value.mobile = 0
        userInfo.value.email = ''
        userInfo.value.role = 0
        userInfo.value.gender = 0
        token.value = ''
    }

    return {
        userInfo,
        token,
        setUserInfo,
        clearUserInfo,
        setToken
    }
}, {
    persist: true // 持久化
})