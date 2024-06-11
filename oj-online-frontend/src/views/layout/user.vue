<template>
  <div class="user">
    <div class="user-header">
      <el-avatar shape="square" :size="size" :src="squareUrl"></el-avatar>
      <div class="user-data">
        <div class="user-name"><strong>用户名</strong></div>
        <div class="user-info">电话</div>
      </div>
    </div>

    <el-divider></el-divider>

    <div class="user-body">
      <div class="user-info-detail">
        <el-descriptions title="个人信息" column=1>
          <el-descriptions-item label="用户名"
            >黑手双城</el-descriptions-item
          >
          <el-descriptions-item label="手机号"
            >18000000000</el-descriptions-item
          >
          <el-descriptions-item label="居住地">成都市</el-descriptions-item>
          <el-descriptions-item label="联系地址"
            >xxxxx大道 1188 号</el-descriptions-item
          >
        </el-descriptions>
      </div>

      <el-divider direction="vertical"></el-divider>
      <div class="problem-solve-info">
        <div>统计信息</div>

        <el-divider></el-divider>
        <div>提交记录</div>
      </div>
    </div>
  </div>
</template>

<script>
import * as echarts from 'echarts'
export default {
  name: 'UserPage',
  data () {
    return {
      size: 80,
      squareUrl:
        'https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png'
    }
  }
}
</script>

<style>
.user-header {
  display: flex; /* 将父元素设置为flex容器 */
}
.user-header .user-data {
  margin-top: 10px;
  margin-left: 20px;
}

.user-body {
  display: flex;
}
.user-body .user-info-detail{
  margin: 30px;
  width: 250px;
}
.user-body .problem-solve-info{
  margin: 30px;
}
</style>
