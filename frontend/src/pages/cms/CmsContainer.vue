<script lang="ts" setup>
import {
    Management,
    UserFilled,
    Flag,
    List,
    BellFilled
} from '@element-plus/icons-vue'
</script>

<template>
    <!--
        后台管理容器：嵌套在主布局（LayoutContainer）内部，
        保留顶部导航栏，此处只提供左侧管理菜单 + 内容区
    -->
    <el-container class="admin-container">
        <el-aside width="200px">
            <el-menu :default-active="$route.path" router class="admin-menu">
                <el-menu-item index="/admin/user">
                    <el-icon>
                        <UserFilled />
                    </el-icon>
                    <span>用户管理</span>
                </el-menu-item>
                <el-menu-item index="/admin/problem">
                    <el-icon>
                        <Management />
                    </el-icon>
                    <span>题目管理</span>
                </el-menu-item>
                <el-menu-item index="/admin/contest">
                    <el-icon>
                        <Flag />
                    </el-icon>
                    <span>比赛管理</span>
                </el-menu-item>
                <el-menu-item index="/admin/template">
                    <el-icon>
                        <List />
                    </el-icon>
                    <span>模板管理</span>
                </el-menu-item>
                <el-menu-item index="/admin/notice">
                    <el-icon>
                        <BellFilled />
                    </el-icon>
                    <span>公告管理</span>
                </el-menu-item>
            </el-menu>
        </el-aside>

        <el-main class="admin-main">
            <!-- 二级路由出口 -->
            <router-view></router-view>
        </el-main>
    </el-container>
</template>

<style lang="less" scoped>
.admin-container {
    height: 100%;
    background-color: #fff;
    border-radius: 8px;
    overflow: hidden;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);

    .el-aside {
        background-color: #fff;
    }

    .el-menu {
        border-right: none;
        background-color: transparent;
    }

    .el-main {
        background-color: #fff;
        padding: 20px;
        overflow: hidden;
    }
}
</style>
