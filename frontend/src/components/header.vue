<template>
  <div>
    <div class="header-container">
      <el-dropdown
        v-if="isLogin"
        placement="bottom-end"
        trigger="click"
        @command="handleCommand"
      >
        <span class="el-dropdown-link">
          {{ dropdownLabelName
          }}<i class="el-icon-arrow-down el-icon--right"></i>
        </span>
        <el-dropdown-menu slot="dropdown">
          <el-dropdown-item command="Home">Home</el-dropdown-item>
          <el-dropdown-item command="Submissions">Submissions</el-dropdown-item>
          <el-dropdown-item command="Settings">Settings</el-dropdown-item>
          <el-dropdown-item command="Logout">Logout</el-dropdown-item>
        </el-dropdown-menu>
      </el-dropdown>
      <el-button v-else @click="handleBtnClick">Login</el-button>
    </div>

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
      dialogVisible: false,
      dialogWidth: '%100'
    }
  },
  computed: {
    isLogin () {
      return this.$store.getters.token
    },
    dropdownLabelName () {
      if (this.isLogin) {
        return this.$store.getters.nickname
      } else {
        return 'Login'
      }
    }
  },
  methods: {
    handleBtnClick () {
      // 判断是否登录
      if (!this.isLogin) {
        this.dialogVisible = true
        return
      }
      this.$router.push('/user')
    },
    handleCommand (command) {
      this.$message('click on item ' + command)
      switch (command) {
        case 'Home':
          this.$router.push('/user')
          break
        case 'Submissions':
          this.$router.push('/submission')
          break
        case 'Settings':
          this.$router.push('/setting')
          break
        case 'Logout':
          this.$store.dispatch('logout')
          this.$router.push('/')
          break
        default:
          break
      }
    }
  }
}
</script>

<style scoped>
.header-container {
  margin-right: 18px;
  cursor: pointer;
  color: #fff;
  font-family: inherit;
  font-size: 16px;
}
.header-container .el-dropdown-link {
  color: #fff;
}
.header-container .el-icon-arrow-down {
  font-size: 15px;
}
.header-container .el-button {
  color: #fff;
  background-color: #158fbf;
  border: none;
  font-size: 18px;
}
.header-container .el-button:hover {
  color: #fffc59;
  background-color: #158fbf;
  border: none;
}
/* 按下样式 */
.header-container .el-button:active {
  color: #ffe44e91;
  background-color: #158fbf;
  border: none;
}
</style>
