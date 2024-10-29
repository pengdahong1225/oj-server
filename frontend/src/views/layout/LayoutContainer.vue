<script setup lang="ts">
import { useRouter } from 'vue-router'
import {
  User,
  Crop,
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
  }
}

</script>

<template>
  <el-container class="layout-container">
    <!-- 头部菜单栏 -->
    <el-header class="top">
      <span class="logo">PGOJ</span>
      <el-menu mode="horizontal" :default-active="$route.path" router>
        <el-menu-item index="/home">Home</el-menu-item>
        <el-menu-item index="/problems">Problems</el-menu-item>
        <el-menu-item index="/contest">Contests</el-menu-item>
        <el-menu-item index="/status">Status</el-menu-item>
      </el-menu>

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
  background-color: #fff;

  .top {
    display: flex;

    .logo {
      font-size: 24px;
      font-weight: bold;
      color: #333;
      padding: 0 24px;
      line-height: 60px;
    }

    .el-menu {
      flex: 1;
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

  .el-footer {
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 14px;
    color: #666;
  }
}
</style>
