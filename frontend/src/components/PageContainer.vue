<script setup lang="ts">

defineProps({
    title: {
        required: true,
        type: String
    }
})
</script>

<template>
    <el-card class="page-container">
        <template #header>
            <div class="header">
                <span><strong>{{ title }}</strong></span>
                <div class="extra">
                    <slot name="extra"></slot>
                </div>
            </div>
        </template>

        <slot></slot>
    </el-card>
</template>

<style lang="less" scoped>
.page-container {
    min-height: 100%;
    box-sizing: border-box;

    .header {
        display: flex;
        justify-content: space-between;
        align-items: center;
    }
}
</style>