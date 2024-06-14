import Vue from 'vue'
import Vuex from 'vuex'

import user from './modules/user'
import problem from './modules/problem'

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
  },
  getters: {
    token: state => {
      return state.user.token
    },
    problemList: state => {
      return state.problem.problemList
    },
    uid: state => {
      return state.user.userInfo.uid
    },
    nickname: state => {
      return state.user.userInfo.nickname
    },
    userSolvedList: state => {
      return state.user.userSolvedList
    }
  },
  mutations: {
  },
  actions: {
  },
  modules: {
    user,
    problem
  }
})
