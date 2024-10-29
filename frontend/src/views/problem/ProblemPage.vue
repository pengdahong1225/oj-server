<script lang="ts" setup>
import { useRoute } from 'vue-router'

const route = useRoute()
console.log(route.params)
</script>

<template>
    <div class="container">
        <!-- 左边题目描述和编辑区域 -->
        <div class="left">
            <el-card class="problem-description" shadow="hover">
                <template #header>
                    <strong>A+B</strong>
                </template>

                <div class="title">
                    Description
                </div>
                <div class="content">
                    输入一个正整数
                </div>
            </el-card>
        </div>

        <!-- 右边题目information区域 -->
        <div class="right">
            <el-card shadow="hover">
                <template #header>
                    <strong>Information</strong>
                </template>
                <el-descriptions :column="1" border>
                    <el-descriptions-item label="ID">1</el-descriptions-item>
                    <el-descriptions-item label="Time Limit">18100000000</el-descriptions-item>
                    <el-descriptions-item label="Memory Limit">Suzhou</el-descriptions-item>
                    <el-descriptions-item label="IO Mode">Suzhou</el-descriptions-item>
                    <el-descriptions-item label="Created By">Suzhou</el-descriptions-item>
                    <el-descriptions-item label="Level">Suzhou</el-descriptions-item>
                    <el-descriptions-item label="Tags">
                        <el-tag size="small">School</el-tag>
                    </el-descriptions-item>
                </el-descriptions>
            </el-card>
        </div>
    </div>
</template>

<style lang="scss" scoped>
.container {
    display: flex;
    width: 80%;
    margin: auto;
    align-items: flex-start;

    /* 确保卡片顶部对齐，而不会拉伸高度 */
    .left {
        width: 85%;
        border: 1px solid #e42518;
        margin-right: 5px;

        .title {
            font-size: 21px;
            font-weight: 500;
            line-height: 30px;
            padding: 5px 10px;
        }
    }

    .right {
        width: 15%;
        margin-left: 5px;
    }
}
</style>