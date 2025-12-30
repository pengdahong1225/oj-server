<script setup lang="ts">
import { useRouter } from 'vue-router'
import {
  User,
  Position,
  EditPen,
  SwitchButton,
  CaretBottom
} from '@element-plus/icons-vue'
import avatar from '@/assets/default.png'
import { useUserStore } from '@/stores'

const router = useRouter()
const userStore = useUserStore()
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
      ).then(() => {
        userStore.clearUserInfo()
        router.push('/login')
      }).catch(() => { })
      break
    case 'admin':
      router.push('/admin')
      break
  }
}

</script>

<template>
  <el-container class="layout-container">
    <!-- 头部菜单栏 -->
    <el-header class="top">
      <!-- logo -->
      <span class="logo">PGOJ</span>

      <!-- 菜单栏 -->
      <el-menu 
        mode="horizontal" 
        :default-active="$route.path" 
        router>
        <el-menu-item index="/home">Home</el-menu-item>
        <el-menu-item index="/problems">Problems</el-menu-item>
        <el-menu-item index="/contest">Contests</el-menu-item>
        <el-menu-item index="/status">Status</el-menu-item>
      </el-menu>

      <!-- 下拉 -->
      <el-dropdown placement="bottom-end" @command="handleCommand">
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
            <el-dropdown-item v-if="userStore.userInfo.role === 1" command="admin" :icon="Position">To CMS</el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </el-header>

    <el-main>
      <!-- 二级路由出口 -->
      <router-view></router-view>
    </el-main>

    <el-footer>PGOJ ©2024 Created by Peter</el-footer>
  </el-container>
</template>

<style lang="less" scoped>
.layout-container {
  height: 100vh;
  background-color: #f5f5f5;

  .top {
    display: flex;
    background-color: #e8e8e8;
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
      background: linear-gradient(to bottom, rgba(232, 232, 232, 0.8), rgba(245, 245, 245, 0));
      pointer-events: none;
    }

    .logo {
      font-size: 24px;
      font-weight: bold;
      padding: 0 24px;
      line-height: 60px;
      background: linear-gradient(135deg, #000 0%, #ffd700 100%);
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
      background-clip: text;
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

// 深度选择器，覆盖Element Plus菜单的选中状态颜色为黑色
:deep(.el-menu--horizontal) {
  .el-menu-item {
    &.is-active {
      color: #000 !important;
      border-bottom-color: #000 !important;
    }
  }
}
</style>
