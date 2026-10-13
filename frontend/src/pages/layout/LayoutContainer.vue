<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import {
  User,
  EditPen,
  SwitchButton,
  CaretBottom,
  HomeFilled,
  Grid,
  Tickets,
  Trophy,
  Promotion,
  DataAnalysis,
  Comment,
  Management,
  Lock
} from '@element-plus/icons-vue'
import avatar from '@/assets/default.png'
import logo from '@/assets/logo.a0924d7d.png'
import { useUserStore } from '@/stores'
import { userLoginService, userLogoutService } from '@/api/userController'
import { loginDialogVisible } from '@/utils/loginDialog'

const router = useRouter()
const userStore = useUserStore()
// 当前高亮菜单：/admin/* 统一高亮到「管理」
const activeMenu = computed(() => {
  const path = router.currentRoute.value.path
  return path.startsWith('/admin') ? '/admin' : path
})
const handleCommand = (command: string) => {
  switch (command) {
    case 'profile':
      router.push(`/user/${command}`)
      break
    case 'password':
      console.log('重置密码')
      break
    case 'logout':
      // 退出
      ElMessageBox.confirm(
        '你确认要进行退出么',
        '温馨提示',
        {
          type: 'warning',
          confirmButtonText: 'OK',
          cancelButtonText: 'Cancel',
        }
      ).then(async () => {
        try {
          await userLogoutService()
        } catch (e) {
          // 后端清除 cookie 失败也继续退出
        }
        userStore.clearUserInfo()
        ElMessage.success('已退出登录')
      }).catch(() => { })
      break
  }
}

// 登录对话框
const loginLoading = ref(false)
const loginForm = ref<API.LoginForm>({
  mobile: '',
  password: ''
})

const handleLogin = async () => {
  if (!loginForm.value.mobile || !loginForm.value.password) {
    ElMessage.warning('请输入用户名和密码')
    return
  }
  loginLoading.value = true
  try {
    const resp = await userLoginService(loginForm.value)
    ElMessage.success('登录成功')
    const info = resp.data.data.user_info
    const user_info = <API.UserInfo>{
      uid: info.uid,
      mobile: info.mobile,
      nickname: info.nick_name,
      gender: info.gender,
      email: info.email,
      avatarUrl: info.avatar_url,
      role: info.role
    }
    userStore.setUserInfo(user_info)
    userStore.setToken(resp.data.data.access_token)
    loginDialogVisible.value = false
    loginForm.value.password = ''
  } catch (e) {
    // 错误提示已在 request 拦截器里处理
  } finally {
    loginLoading.value = false
  }
}

const handleForgotPassword = () => {
  ElMessage.info('忘记密码功能暂未实现')
}

</script>

<template>
  <el-container class="layout-container">
    <!-- 头部菜单栏 -->
    <el-header class="top">
      <!-- logo -->
      <img class="logo" :src="logo" alt="OpenOJ" />

      <!-- 菜单栏 -->
      <el-menu
        mode="horizontal"
        :default-active="activeMenu"
        router>
        <el-menu-item index="/home">
          <el-icon><HomeFilled /></el-icon>
          <span>首页</span>
        </el-menu-item>
        <el-menu-item index="/problems">
          <el-icon><Grid /></el-icon>
          <span>题目</span>
        </el-menu-item>
        <el-menu-item index="/training">
          <el-icon><Tickets /></el-icon>
          <span>训练</span>
        </el-menu-item>
        <el-menu-item index="/contest">
          <el-icon><Trophy /></el-icon>
          <span>比赛</span>
        </el-menu-item>
        <el-menu-item index="/judge">
          <el-icon><Promotion /></el-icon>
          <span>测评</span>
        </el-menu-item>
        <el-menu-item index="/rank">
          <el-icon><DataAnalysis /></el-icon>
          <span>排名</span>
        </el-menu-item>
        <el-menu-item index="/discussion">
          <el-icon><Comment /></el-icon>
          <span>讨论</span>
        </el-menu-item>
        <el-menu-item v-if="userStore.userInfo.role === 1" index="/admin">
          <el-icon><Management /></el-icon>
          <span>管理</span>
        </el-menu-item>
      </el-menu>

      <!-- 未登录：登录按钮 -->
      <el-button v-if="!userStore.token" class="login-btn" type="primary" @click="loginDialogVisible = true">登录</el-button>

      <!-- 已登录：下拉菜单 -->
      <el-dropdown v-else placement="bottom-end" @command="handleCommand">
        <span class="el-dropdown__box">
          <el-avatar :src="avatar" />
          <el-icon>
            <CaretBottom />
          </el-icon>
        </span>
        <!-- 折叠的下拉部分 -->
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item command="profile" :icon="User">基本资料</el-dropdown-item>
            <el-dropdown-item command="password" :icon="EditPen">重置密码</el-dropdown-item>
            <el-dropdown-item command="logout" :icon="SwitchButton">退出登录</el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </el-header>

    <el-main>
      <!-- 二级路由出口 -->
      <router-view></router-view>
    </el-main>

    <el-footer>OpenOJ ©2024 Created by Peter</el-footer>
  </el-container>

  <!-- 登录对话框 -->
  <el-dialog v-model="loginDialogVisible" width="400px" :close-on-click-modal="false" class="login-dialog">
    <template #header>
      <div class="dialog-title">登录 - OpenOJ</div>
    </template>
    <div class="login-body">
      <el-input v-model="loginForm.mobile" placeholder="用户名" :prefix-icon="User" size="large" clearable />
      <el-input v-model="loginForm.password" type="password" placeholder="密码" :prefix-icon="Lock" size="large" show-password
        clearable @keyup.enter="handleLogin" />
      <el-button type="primary" :loading="loginLoading" size="large" class="login-submit" @click="handleLogin">登录</el-button>
      <el-button link type="primary" class="forgot-btn" @click="handleForgotPassword">忘记密码</el-button>
    </div>
  </el-dialog>
</template>

<style lang="less" scoped>
.layout-container {
  height: 100vh;
  background-color: #f5f5f5;

  .top {
    display: flex;
    background-color: #fff;
    position: relative;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
    z-index: 1;

    // 添加渐变效果
    &::after {
      content: '';
      position: absolute;
      bottom: -10px;
      left: 0;
      right: 0;
      height: 10px;
      background: linear-gradient(to bottom, rgba(0, 0, 0, 0.05), rgba(245, 245, 245, 0));
      pointer-events: none;
    }

    .logo {
      height: 40px;
      margin: 0 24px;
      align-self: center;
    }

    .login-btn {
      align-self: center;
      margin-right: 24px;
    }

    .el-menu {
      flex: 1;
      background-color: transparent;
    }

    .el-dropdown {
      .el-dropdown__box {
        display: flex;
        align-items: center;

        .el-icon {
          color: #999;
          margin-left: 10px;
        }

        &:active,
        &:focus {
          outline: none;
        }
      }
    }
  }

  .el-main {
    background-color: #f5f5f5;
    min-height: 0;
  }

  .el-footer {
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 14px;
    color: #666;
  }
}

:deep(.el-menu--horizontal) {
  .el-menu-item {
    &.is-active {
      color: var(--el-color-primary) !important;
      border-bottom-color: var(--el-color-primary) !important;
    }
  }
}

.dialog-title {
  text-align: center;
  font-weight: bold;
  font-size: 20px;
}

.login-body {
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin-top: 24px;

  // 输入框边框颜色
  --el-input-border-color: rgb(14 116 144);
  --el-input-hover-border-color: rgb(14 116 144);
  --el-input-focus-border-color: rgb(14 116 144);
}

.login-submit {
  width: 100%;
}

.forgot-btn {
  align-self: flex-end;
}
</style>

<style lang="less">
.login-dialog {
  border-radius: 12px !important;
}
</style>
