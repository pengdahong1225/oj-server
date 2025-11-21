<script lang="ts" setup>
import { ref, onMounted } from "vue"
import { Search, Select } from '@element-plus/icons-vue'
import { queryProblemListService, getProblemTagListService } from '@/api/problemController'
import { useUserStore } from '@/stores'

const userStore = useUserStore()
const loading = ref(false)
const total = ref(0)
const problemList = ref<API.Problem[]>([])
const tag_list = ref<string[]>([])
const params = ref(<API.QueryProblemListParams>{
    page: 1,
    page_size: 10, // page_size默认为10
    keyword: '',
    tag: ''
})
onMounted(() => {
    getProblemList()
    getProblemTagList()
})
const getProblemList = async () => {
    if (userStore.userInfo.uid > 0) {
        params.value.uid = userStore.userInfo.uid
    }
    loading.value = true
    const res = await queryProblemListService(params.value)
    total.value = res.data.data.total
    problemList.value = res.data.data.list
    loading.value = false
}
const getProblemTagList = async () => {
    const res = await getProblemTagListService()
    if (Array.isArray(res.data.data) && res.data.data.length > 0) {
        tag_list.value = res.data.data
    }
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
const onTagClick = (tag: string) => {
    params.value.tag = tag
    getProblemList()
}
</script>

<template>
    <div class="container">
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
                <el-table-column label="Status" prop="status" width="80" align="center">
                    <template #default="{ row }">
                        <el-icon v-if="row.status === 1" color="green" size="18"><Select /></el-icon>
                    </template>
                </el-table-column>

                <el-table-column label="#" prop="problem_id" width="100">
                    <template #default="{ row }">
                        <el-link type="primary" :underline="false" @click="
                            $router.push({
                                path: `/problem/${row.problem_id}`
                            })
                            ">{{ row.problem_id }}</el-link>
                    </template>
                </el-table-column>

                <el-table-column label="Title" prop="title" width="180">
                    <template #default="{ row }">
                        <el-link type="primary" :underline="false" @click="
                            $router.push({
                                path: `/problem/${row.problem_id}`
                            })
                            ">{{ row.problem_title }}</el-link>
                    </template>
                </el-table-column>
                
                <el-table-column label="Tags" width="500">
                    <template #default="{ row }">
                        <el-tag v-for="tag in row.tags" :key="tag" style="margin-left: 3px;margin-right: 3px;">{{ tag }}</el-tag>
                    </template>
                </el-table-column>

                <el-table-column label="Level" prop="level" width="200">
                    <template #default="{ row }">
                        <el-tag v-if="row.level === 1" type="primary">简单</el-tag>
                        <el-tag v-else-if="row.level === 2" type="warning">中等</el-tag>
                        <el-tag v-else type="danger">困难</el-tag>
                    </template>
                </el-table-column>

                <el-table-column label="通过率">
                   
                </el-table-column>

                <template #empty>
                    <el-empty description="没有数据"></el-empty>
                </template>
            </el-table>

            <!-- 分页 -->
            <el-pagination v-model:current-page="params.page" v-model:page-size="params.page_size" :total="total"
                :background="true" layout="prev, pager, next, jumper" @current-change="handleCurrentChange"
                style="margin-top: 20px; justify-content: flex-end;" />
        </el-card>
        <!-- 右侧tag池区域 -->
        <el-card class="tag-pool" shadow="hover" header="Tag Pool">
            <el-button v-for="item in tag_list" :key="item" type="info" plain round @click="onTagClick(item)">
                {{ item }}</el-button>
        </el-card>
    </div>
</template>

<style lang="less" scoped>
.container {
    display: flex;
    width: 100%;
    align-items: flex-start;
    /* 确保卡片顶部对齐，而不会拉伸高度 */

    .problem-list {
        width: 75%;
        /* 占 70% 的宽度 */
        margin-right: 6px;
    }

    .tag-pool {
        width: 25%;
        /* 占 30% 的宽度 */
        margin-left: 6px;

        .el-button {
            margin: 5px;
        }
    }
}

.form {
    margin-left: 10px;

    .el-form-item {
        margin-bottom: 0;
    }
}
</style>