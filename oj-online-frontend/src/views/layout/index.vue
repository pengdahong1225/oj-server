<template>
  <div>
    <div class="top">
      <!-- 选项栏 -->
      <el-menu
        active-text-color="#dfff7d"
        background-color="#158fbf"
        class="el-menu-vertical"
        :default-active="activeIndex"
        text-color="#fff"
        router
        mode="horizontal"
      >
        <el-menu-item index="0" disabled>
          <img style="width: 50px" src="@/assets/Cup.png" alt="Oj Online" />
        </el-menu-item>
        <div class="flex-grow" />
        <el-menu-item index="1" route="/home">主页</el-menu-item>
        <el-menu-item index="2" route="/problemset">题目列表</el-menu-item>
        <el-menu-item index="3" route="/rank">排行榜</el-menu-item>
        <!-- <el-menu-item index="4" route="/login">登录/注册</el-menu-item> -->
      </el-menu>
      <!-- User -->
      <Header class="topHeader"></Header>
    </div>

    <el-divider></el-divider>

    <!-- 二级路由出口：二级组件展示的位置 -->
    <router-view></router-view>
  </div>
</template>

<script>
import Header from '@/components/header.vue'
export default {
  name: 'LayoutIndex',
  components: {
    Header
  },
  data () {
    return {
      activeIndex: '1'
    }
  },
  methods: {
    handleSelect (key, keyPath) {
      // index: 选中菜单项的 index, indexPath: 选中菜单项的 index path
      console.log(key, keyPath)
    }
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
</style>
