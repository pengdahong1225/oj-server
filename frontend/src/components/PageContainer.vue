<script setup lang="ts">

defineProps({
    title: {
        required: true,
        type: String
    }
})
</script>

<template>
    <div class="page-container">
        <div class="header">
            <span><strong>{{ title }}</strong></span>
            <div class="extra">
                <slot name="extra"></slot>
            </div>
        </div>
        <div class="body">
            <slot></slot>
        </div>
    </div>
</template>

<style lang="less" scoped>
.page-container {
    height: 100%;
    display: flex;
    flex-direction: column;

    .header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding-bottom: 16px;
        flex-shrink: 0;
    }

    .body {
        flex: 1;
        min-height: 0;
        overflow: auto;
    }
}
</style>
