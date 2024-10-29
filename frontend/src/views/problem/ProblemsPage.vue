<script lang="ts" setup>
import { ref, onMounted } from "vue"
import { Search, Select } from '@element-plus/icons-vue'
import type { Problem, QueryProblemListParams } from '@/types/problem'
import type { UPSSParams } from '@/types/user'
import { queryProblemListService } from '@/api/problem'
import { upssService } from '@/api/user'

const loading = ref(false)
const total = ref(0)
const problemList = ref<Problem[]>([
    {
        id: 1,
        title: 'Hello World',
        description: 'Hello World',
        level: 1,
        tags: ['Java', 'Python'],
        status: 1
    },
    {
        id: 2,
        title: 'tow sum',
        description: 'tow sum',
        level: 2,
        tags: ['Java', 'Python'],
        status: 0,
    },
    {
        id: 3,
        title: '接雨水',
        description: '接雨水',
        level: 3,
        tags: ['Java', 'Python'],
        status: 0
    }
])
const tag_list = ref([
    'Java', 'Python', 'C++', 'JavaScript', 'TypeScript', 'Go', 'Rust', 'Swift', 'Kotlin', 'PHP', 'Ruby', 'C#', 'SQL', 'HTML', 'CSS', 'Dart', 'Objective-C', 'R', 'Matlab', 'Lua', 'Vue', 'React', 'Angular', 'Node.js', 'Express', 'Flask', 'Django', 'Spring', 'Hibernate'
])
const params = ref(<QueryProblemListParams>{
    page: 1,
    page_size: 10, // page_size默认为10
    keyword: '',
    tag: ''
})
onMounted(() => {
    // getProblemList()
    // getUpss()
})
const getProblemList = async () => {
    loading.value = true
    const res = await queryProblemListService(params.value)
    console.log(res)
    total.value = res.data.total
    problemList.value = res.data.data
    loading.value = false
}
const getUpss = async () => {
    if (problemList.value.length === 0) {
        return
    }

    let params = <UPSSParams>({
        uid: 1
    })
    problemList.value.forEach((item) => {
        params.problem_ids.push(item.id)
    })

    const res = await upssService(params)
    console.log(res)
}
const handleCurrentChange = (page: number) => {
    params.value.page = page
    getProblemList()
    getUpss()
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
    console.log(tag)
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
                <el-table-column label="Status" width="80" align="center">
                    <template #default="{ row }">
                        <el-icon v-if="row.status === 1" color="green" size="18"><Select /></el-icon>
                    </template>
                </el-table-column>

                <el-table-column label="#" prop="id" width="80">
                    <template #default="{ row }">
                        <el-link type="primary" :underline="false" @click="
                            $router.push({
                                path: `/problem/${row.id}`
                            })
                            ">{{ row.id }}</el-link>
                    </template>
                </el-table-column>

                <el-table-column label="Title" prop="title">
                    <template #default="{ row }">
                        <el-link type="primary" :underline="false" @click="
                            $router.push({
                                path: `/problem/${row.id}`
                            })
                            ">{{ row.title }}</el-link>
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
            <el-button v-for="item in tag_list" :key="item" type="info" plain round="true" @click="onTagClick(item)">
                {{ item }}</el-button>
        </el-card>
    </div>
</template>

<style lang="scss" scoped>
.container {
    display: flex;
    width: 100%;
    align-items: flex-start; /* 确保卡片顶部对齐，而不会拉伸高度 */

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