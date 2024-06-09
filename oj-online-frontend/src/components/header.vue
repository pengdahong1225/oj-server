<template>
  <div>
    <!-- 登录了就显示头像，没登录就显示默认头像点击提示登录 -->
    <el-avatar
      :shape="shape"
      :size="size"
      :src="avatarUrl"
      v-on:click.native="handleAvatarClick"
    >
      <img
        src="https://cube.elemecdn.com/e/fd/0fc7d20532fdaf769a25683617711png.png"
      />
    </el-avatar>

    <el-dialog
      title="Welcome OJ Onine"
      :visible.sync="dialogVisible"
      :width="dialogWidth"
    >
      <Login @close="dialogVisible = false"></Login>
    </el-dialog>
  </div>
</template>

<script>
import Login from '@/components/login.vue'
export default {
  name: 'HeaderComponent',
  components: {
    Login
  },
  data () {
    return {
      avatarUrl:
        'https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png',
      shape: 'circle',
      size: 50,
      dialogVisible: false,
      dialogWidth: '%100'
    }
  },
  computed: {
    isLogin () {
      return this.$store.getters.token
    }
  },
  methods: {
    handleAvatarClick () {
      // 判断是否登录
      if (!this.isLogin) {
        this.dialogVisible = true
        return
      }
      this.$router.push('/user')
    }
  }
}
</script>

<style scoped>
.el-avatar {
  cursor: pointer;
  margin-right: 20px;
}
</style>
