<script lang="ts" setup>

</script>

<template>
    <div class="container">
        <div class="d-title">
            Description
        </div>
        <div class="d-content">
            输入一个正整数
        </div>
    </div>
</template>

<style lang="scss" scoped>
.container {
    background: rgb(251, 251, 251);
    width: 100%;
    display: flex;
    flex-direction: column; /* 竖直排列 */
    align-items: flex-start; /* 左对齐 */

    .d-title {
        font-size: 21px;
        font-weight: 500;
        line-height: 30px;
        padding: 5px 10px;
        color: #3091f2;
    }
    .d-content {
      margin-left: 25px;
      margin-right: 20px;
    }
}
</style>