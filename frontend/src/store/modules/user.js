import { setUserInfo, setToken, getUserInfo, getToken } from '@/utils/storage'

export default {
  namespaced: true,
  state () {
    return {
      token: getToken(),
      userInfo: getUserInfo(),
      userSolvedList: []
    }
  },
  mutations: {
    // 所有mutations的第一个参数，都是state
    setUserInfo (state, obj) {
      state.userInfo = obj
      setUserInfo(obj)
    },
    setToken (state, obj) {
      state.token = obj
      setToken(obj)
    },
    setUserSolvedList (state, obj) {
      state.userSolvedList = obj
    }
  },
  actions: {
    logout (context) {
      // 重置
      context.commit('setUserInfo', {})
      context.commit('setToken', '')
    }
  },
  getters: {}
}
