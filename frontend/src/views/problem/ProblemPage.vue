<script lang="ts" setup>
import { onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { queryProblemDetailService } from '@/api/problem'
import type { Problem } from '@/types/problem'

onMounted(() => {
    getProblemDetail()
})

const route = useRoute()
const problem = ref<Problem>({
    id: 0,
    title: '',
    description: '',
    level: 0,
    tags: [],
    status: 1
})
const getProblemDetail = async () => {
    const res = await queryProblemDetailService(Number(route.params.id))
    problem.value.id = res.data.data.id
    problem.value.title = res.data.data.title
    problem.value.description = res.data.data.description
    problem.value.level = res.data.data.level
    problem.value.tags = res.data.data.tags
    problem.value.status = res.data.data.status
}

</script>

<template>
    <div class="container">
        <!-- 左边题目描述和编辑区域 -->
        <div class="left">
            <el-card class="problem-description" shadow="hover">
                <template #header>
                    <strong>A+B</strong>
                </template>

                <descriptionItem style="width: 100%;"></descriptionItem>

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
        width: 80%;
        border: 1px solid #e42518;
        margin-right: 5px;
    }

    .right {
        width: 20%;
        margin-left: 5px;
    }
}
</style>