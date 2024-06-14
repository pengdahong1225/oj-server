<template>
  <div>
    <div class="top">
      <!-- 选项栏 -->
      <el-menu
        active-text-color="#dfff7d"
        background-color="#158fbf"
        class="el-menu-vertical"
        :default-active="$route.path"
        text-color="#fff"
        router
        mode="horizontal"
      >
        <el-menu-item index="0" disabled>
          <img style="width: 50px" src="@/assets/Cup.png" alt="Oj Online" />
        </el-menu-item>
        <div class="flex-grow" />
        <el-menu-item index="/home" route="/home">主页</el-menu-item>
        <el-menu-item index="/problemlist" route="/problemlist">题目列表</el-menu-item>
        <el-menu-item index="/rank" route="/rank">排行榜</el-menu-item>
      </el-menu>
      <!-- User -->
      <Header class="topHeader"></Header>
    </div>

    <el-divider></el-divider>

    <div class="container">
      <!-- 二级路由出口：二级组件展示的位置 -->
      <router-view></router-view>
    </div>
  </div>
</template>

<script>
import Header from '@/components/header.vue'
export default {
  name: 'LayoutIndex',
  components: {
    Header
  }
}
</script>

<style>
.top {
  /* 弹性盒子 */
  display: flex;
  align-items: center;
  background-color: #158fbf;
}
.topHeader {
  flex-grow: 50px;
}
.el-menu-vertical {
  flex: 1;
}
.container {
  width: 90%; /* 根据需要设置宽度 */
  margin: 0 auto; /* 上下边距为0，左右边距自动，使得.main盒子居中 */
  background-color: #fff; /* 设置背景颜色 */
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.1); /* 设置阴影效果 */
}
</style>
