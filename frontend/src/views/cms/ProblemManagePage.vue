<script setup lang="ts">
import { Delete, Edit, Search } from '@element-plus/icons-vue'
import { onMounted, ref, watch } from 'vue'
import { queryProblemListService } from '@/api/problem'

import type { Problem, QueryProblemListParams } from '@/types/problem'
import { formatTime } from '@/utils/format'
import ProblemEdit from '@/components/problemEdit.vue'

onMounted(() => {
    getProblemList()
})

const loading = ref(false)
// 总数
const total = ref(0)
// 文章列表数据
const problemList = ref<Problem[]>([])
// 分页请求参数
const params = ref(<QueryProblemListParams>{
    page: 1,
    page_size: 10, // page_size默认为10
    keyword: '',
    tag: ''
})

const getProblemList = async () => {
    loading.value = true
    const res = await queryProblemListService(params.value)
    console.log(res)
    total.value = res.data.data.total
    problemList.value = res.data.data.data
    loading.value = false
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

const problemEditRef = ref()
const onAddProblem = () => {
    problemEditRef.value.open({})
}
const onEdit = (row: Problem) => {
    problemEditRef.value.open(row)
}
const onDelete = (row: Problem) => {
    console.log(row)
}
// 添加或者编辑 成功的回调
const onSuccess = (isAdd: boolean) => {
    if (isAdd) {
        // 如果是添加，最好渲染最后一页
        const lastPage = Math.ceil((total.value + 1) / params.value.page_size)
        // 更新成最大页码数，再渲染
        params.value.page = lastPage
        ElMessage.success('添加成功')
    } else {
        ElMessage.success('更新成功')
    }

    getProblemList()
}

</script>

<template>
    <page-container title="题目列表">
        <template #extra>
            <el-button type="primary" @click="onAddProblem">添加题目</el-button>
        </template>

        <!-- 表单区域 inline属性代表表单元素都在一行显示 -->
        <el-form inline>
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
            <el-table-column label="#" prop="id" width="80">
                <template #default="{ row }">
                    {{ row.id }}
                </template>
            </el-table-column>

            <el-table-column label="Title" prop="title">
                <template #default="{ row }">
                    {{ row.title }}
                </template>
            </el-table-column>
            <el-table-column label="Level" prop="level">
                <template #default="{ row }">
                    <el-tag v-if="row.level === 1" type="primary">简单</el-tag>
                    <el-tag v-else-if="row.level === 2" type="warning">中等</el-tag>
                    <el-tag v-else type="danger">困难</el-tag>
                </template>
            </el-table-column>
            <el-table-column label="Tags">
                <template #default="{ row }">
                    <el-tag v-for="tag in row.tags" :key="tag" style="margin-left: 3px;margin-right: 3px;">{{ tag
                        }}</el-tag>
                </template>
            </el-table-column>
            <el-table-column label="Created At">
                <template #default="{ row }">
                    {{ formatTime(row.create_at) }}
                </template>
            </el-table-column>
            <el-table-column label="Created By">
                <template #default="{ row }">
                    {{ row.create_by }}
                </template>
            </el-table-column>
            <el-table-column label="操作">
                <template #default="{ row }">
                    <el-button circle plain type="primary" :icon="Edit" @click="onEdit(row)"></el-button>
                    <el-button circle plain type="danger" :icon="Delete" @click="onDelete(row)"></el-button>
                </template>
            </el-table-column>

            <template #empty>
                <el-empty description="没有数据"></el-empty>
            </template>
        </el-table>

        <!-- 分页 -->
        <el-pagination v-model:current-page="params.page" v-model:page-size="params.page_size" :total="total"
            :background="true" layout="prev, pager, next, jumper" @current-change="handleCurrentChange"
            style="margin-top: 20px; justify-content: flex-end;" />

        <!-- 抽屉 -->
        <ProblemEdit ref="problemEditRef" @success="onSuccess"></ProblemEdit>
    </page-container>
</template>

<style lang="scss" scoped></style>