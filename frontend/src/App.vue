<script setup lang="ts">
</script>

<template>
  <div>
    <router-view></router-view>
  </div>
</template>

<style>
:root {
  --el-color-primary: #0E7490 !important;
  --el-color-primary-light-3: #569EB1 !important;
  --el-color-primary-light-5: #87BAC8 !important;
  --el-color-primary-light-7: #B7D5DE !important;
  --el-color-primary-light-8: #CFE3E9 !important;
  --el-color-primary-light-9: #E7F1F4 !important;
  --el-color-primary-dark-2: #0B5D73 !important;
}

* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html, body {
  margin: 0;
  padding: 0;
  height: 100%;
  width: 100%;
  overflow-x: hidden;
}

#app {
  width: 100%;
  height: 100%;
}
</style>
