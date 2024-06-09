export default {
  namespaced: true,
  state () {
    return {
      totalSize: 0,
      problemList: [],
      cursor: 1
    }
  },
  mutations: {
    // 所有mutations的第一个参数，都是state
    setProblemInfo (state, obj) {
      state.totalSize = obj.total
      state.problemList = obj.data
      state.cursor = obj.cursor
    }
  },
  actions: {},
  getters: {}
}
