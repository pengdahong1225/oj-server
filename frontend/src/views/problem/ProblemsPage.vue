<script lang="ts" setup>
import { ref, onMounted } from "vue"
import { Search } from '@element-plus/icons-vue'
import type { Problem, QueryProblemListParams } from '@/types/problem'
import { queryProblemListService } from '@/api/problem'

const loading = ref(false)
const total = ref(0)
const problemList = ref<Problem[]>([])
const params = ref(<QueryProblemListParams>{
    page: 1,
    page_size: 10,
    keyword: '',
    tag: ''
})
onMounted(() => {
    // getProblemList()
})
const getProblemList = async () => {
    loading.value = true
    const res = await queryProblemListService(params.value)
    console.log(res)
    total.value = res.data.total
    problemList.value = res.data.data
    loading.value = false
}
const handleSizeChange = (size: number) => {
    // 只要是每页的条数改变了 => 重新从第一页开始渲染
    params.value.page = 1
    params.value.page_size = size
    getProblemList()
}
const handleCurrentChange = (page: number) => {
    params.value.page = page
    getProblemList()
}
const onSearch = () => {
    params.value.page = 1
    getProblemList()
}
const onReset = () => {
    params.value.page = 1
    params.value.keyword = ''
    params.value.tag = ''
    getProblemList()
}
</script>

<template>
    <el-space class="container">
        <!-- 题目列表区域 -->
        <el-card class="problem-list" shadow="hover">
            <!-- 表单区域 -->
            <el-form inline class="form">
                <el-form-item style="margin-right: 15px;">
                    <el-input v-model="params.keyword" style="width: 240px" placeholder="problem name"
                        :suffix-icon="Search" />
                </el-form-item>
                <el-form-item>
                    <el-button @click="onSearch" type="primary">搜索</el-button>
                    <el-button @click="onReset">重置</el-button>
                </el-form-item>
            </el-form>

            <!-- 表格区域 -->
            <el-table v-loading="loading" :data="problemList">
                <el-table-column label="Title" prop="title">
                    <template #default="{ row }">
                        <el-link type="primary" :underline="false">{{ row.title }}</el-link>
                    </template>
                </el-table-column>
                <el-table-column label="Level" prop="level">
                    <template #default="{ row }">
                        <el-tag v-if="row.level === 1" type="danger">简单</el-tag>
                        <el-tag v-else-if="row.level === 2" type="warning">中等</el-tag>
                        <el-tag v-else type="primary">困难</el-tag>
                    </template>
                </el-table-column>
                <el-table-column label="Tags">
                    <template #default="{ row }">
                        <el-tag v-for="tag in row.tags" :key="tag">{{ tag }}</el-tag>
                    </template>
                </el-table-column>

                <template #empty>
                    <el-empty description="没有数据"></el-empty>
                </template>
            </el-table>

            <!-- 分页 -->
            <el-pagination v-model:current-page="params.page" v-model:page-size="params.page_size"
                :page-sizes="[2, 3, 5, 10]" :background="true" layout="total, sizes, prev, pager, next, jumper"
                :total="total" @size-change="handleSizeChange" @current-change="handleCurrentChange"
                style="margin-top: 20px; justify-content: flex-end;" />
        </el-card>
        <!-- 右侧tag池区域 -->
        <el-card class="tag-pool" shadow="hover" header="Tag Pool" style="width: 480px">

        </el-card>
    </el-space>
</template>

<style lang="scss" scoped>
.container{
    display: flex;
    width: 100%;
    .problem-list{
        flex: 7; /* 占 70% 的宽度 */
    }
    .tag-pool{
        flex: 3; /* 占 30% 的宽度 */
    }
}

.form {
    margin-left: 10px;

    .el-form-item {
        margin-bottom: 0;
    }
}
</style>