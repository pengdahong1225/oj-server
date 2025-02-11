<script lang="ts" setup>
import {
    Management,
    UserFilled,
    User,
    Position,
    EditPen,
    SwitchButton,
    CaretBottom,
    Flag,
    List,
    BellFilled
} from '@element-plus/icons-vue'
import avatar from '@/assets/default.png'
import { useUserStore } from '@/stores'
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
const userStore = useUserStore()
const router = useRouter()

onMounted(() => {
    //   userStore.getUser()
})
const handleCommand = (command: string) => {
    switch (command) {
        case 'profile':
            router.push(`/admin/${command}`)
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
        case 'home':
            router.push('/home')
            break
    }
}

</script>

<template>
    <!-- 
    el-menu 整个菜单组件
      :default-active="$route.path"  配置默认高亮的菜单项
      router  router选项开启，el-menu-item 的 index 就是点击跳转的路径

    el-menu-item 菜单项
      index="" 配置的是访问的跳转路径，配合default-active的值，实现高亮
    -->
    <el-container class="layout-container">
        <el-aside width="200px">
            <div class="el-aside__logo"></div>
            <el-menu active-text-color="#ffd04b" background-color="rgb(81, 90, 110)" :default-active="$route.path"
                text-color="#fff" router>
                <el-menu-item index="/admin/user">
                    <el-icon>
                        <UserFilled />
                    </el-icon>
                    <span>User</span>
                </el-menu-item>
                <el-menu-item index="/admin/problem">
                    <el-icon>
                        <Management />
                    </el-icon>
                    <span>Problem</span>
                </el-menu-item>
                <el-menu-item index="/admin/contest">
                    <el-icon>
                        <Flag />
                    </el-icon>
                    <span>Contest</span>
                </el-menu-item>
                <el-menu-item index="/admin/template">
                    <el-icon>
                        <List />
                    </el-icon>
                    <span>Judge Template</span>
                </el-menu-item>
                <el-menu-item index="/admin/notice">
                    <el-icon>
                        <BellFilled />
                    </el-icon>
                    <span>Notice</span>
                </el-menu-item>
            </el-menu>
        </el-aside>

        <el-container>
            <el-header>
                <!-- 用户基本信息 -->
                <div>
                    管理员：<strong>{{
                        userStore.userInfo.nickname || userStore.userInfo.mobile
                    }}</strong>
                </div>
                <el-dropdown placement="bottom-end" @command="handleCommand">
                    <span class="el-dropdown__box">
                        <el-avatar :src="userStore.userInfo.avatar_url || avatar" />
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
                            <el-dropdown-item command="home" :icon="Position">To Home</el-dropdown-item>
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
    </el-container>
</template>

<style lang="less" scoped>
.layout-container {
    height: 100vh;

    .el-aside {
        background-color: rgb(81, 90, 110);

        &__logo {
            height: 120px;
            background: url('@/assets/cat.png') no-repeat center / 120px auto;
        }

        .el-menu {
            border-right: none;
        }
    }

    .el-header {
        background-color: #fff;
        display: flex;
        align-items: center;
        justify-content: space-between;

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

    .el-footer {
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 14px;
        color: #666;
    }
}
</style>
