<template>
  <div class="user">
    <div class="user-header">
      <el-avatar shape="square" :size="size" :src="squareUrl"></el-avatar>
      <div class="user-data">
        <div class="user-name"> 黑手双城 </div>
        <div class="user-info"> 2102533303@qq.com </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'UserPage',
  data () {
    return {
      size: 80,
      squareUrl: 'https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png'
    }
  }
}
</script>

<style>
</style>
