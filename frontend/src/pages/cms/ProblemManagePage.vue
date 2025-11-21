<script setup lang="ts">
import { Delete, Edit, Search, Open, TurnOff, Upload } from '@element-plus/icons-vue'
import { onMounted, ref } from 'vue'

import { queryProblemListService, deleteProblemService, publishProblemService, hideProblemService } from '@/api/problemController'
import { formatTime } from '@/utils/format'
import ProblemEdit from '@/components/problemEdit.vue'
import { useUserStore } from '@/stores'

onMounted(() => {
    getProblemList()
})

const loading = ref(false)
// 总数
const total = ref(0)
// 文章列表数据
const problemList = ref<API.Problem[]>([])
// 分页请求参数
const params = ref(<API.QueryProblemListParams>{
    page: 1,
    page_size: 10, // page_size默认为10
    keyword: '',
    tag: ''
})

const getProblemList = async () => {
    loading.value = true
    const resp = await queryProblemListService(params.value)
    total.value = resp.data.data.total
    problemList.value = resp.data.data.list
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
    problemEditRef.value.open("create", {})
}
const onEdit = (row: API.Problem) => {
    problemEditRef.value.open("update", row.problem_id)
}
const onPublish = async (id: number) => {
    const resp = await publishProblemService(id)
    if (resp.data.code === 0) {
        ElMessage.success("发布成功")
    } else {
        ElMessage.error(resp.data.message || '发布失败')
    }
    getProblemList()
}
const onHide = async (id: number) => {
    const resp = await hideProblemService(id)
    if (resp.data.code === 0) {
        ElMessage.success("隐藏成功")
    } else {
        ElMessage.error(resp.data.message || '隐藏失败')
    }
    getProblemList()
}
const onDelete = async (id: number) => {
    const resp = await deleteProblemService(id)
    if (resp.data.code === 0) {
        ElMessage.success("删除成功")
        getProblemList()
    } else {
        ElMessage.error(resp.data.message || '删除失败')
    }
}
const uploadHeaders = {
    Authorization: 'Bearer ' + useUserStore().token
}
// 上传回调
const onUploadSuccess = (response: any, row: API.Problem) => {
    if (response.code === 0) {
        ElMessage.success(`题目 ${row.problem_id} 配置上传成功`)
    } else {
        ElMessage.error(response.message || '上传失败')
    }
}
// 添加/编辑的回调
const onSuccess = (mode: string) => {
    if (mode === 'create') {
        // 如果是添加，最好渲染最后一页
        const lastPage = Math.ceil((total.value + 1) / params.value.page_size)
        // 更新成最大页码数，再渲染
        params.value.page = lastPage
    }

    getProblemList()
}

</script>

<template>
    <page-container title="题目管理">
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
            <el-table-column label="#" prop="id" width="60">
                <template #default="{ row }">
                    {{ row.problem_id }}
                </template>
            </el-table-column>

            <el-table-column label="Title" prop="title" width="150">
                <template #default="{ row }">
                    {{ row.problem_title }}
                </template>
            </el-table-column>
            <el-table-column label="Level" prop="level" width="100">
                <template #default="{ row }">
                    <el-tag v-if="row.level === 1" type="primary">简单</el-tag>
                    <el-tag v-else-if="row.level === 2" type="warning">中等</el-tag>
                    <el-tag v-else type="danger">困难</el-tag>
                </template>
            </el-table-column>
            <el-table-column label="Tags" width="500">
                <template #default="{ row }">
                    <el-tag v-for="tag in row.tags" :key="tag" style="margin-left: 3px;margin-right: 3px;">{{ tag
                    }}</el-tag>
                </template>
            </el-table-column>
            <el-table-column label="Created At" with="100">
                <template #default="{ row }">
                    {{ formatTime(row.create_at) }}
                </template>
            </el-table-column>
            <el-table-column label="Updated At" with="100">
                <template #default="{ row }">
                    {{ formatTime(row.update_at) }}
                </template>
            </el-table-column>
            <el-table-column label="操作">
                <template #default="{ row }">
                    <!-- 编辑 -->
                    <el-button circle plain type="primary" :icon="Edit" @click="onEdit(row)"></el-button>
                    <!-- 发布 -->
                    <el-button v-if="row.status === 1" circle plain type="success" :icon="Open"
                        @click="onHide(row.problem_id)"></el-button>
                    <!-- 隐藏 -->
                    <el-button v-else circle plain type="info" :icon="TurnOff"
                        @click="onPublish(row.problem_id)"></el-button>
                    <!-- 删除 -->
                    <el-button circle plain type="danger" :icon="Delete" @click="onDelete(row.problem_id)"></el-button>
                    <!-- 上传配置 -->
                    <el-upload :action="`http://localhost:9000/api/v1/problem/upload_config`" :headers="uploadHeaders"
                        :data="{problem_id: row.problem_id}"
                        :name="'config_file'"
                        :show-file-list="false" :on-success="(resp) => onUploadSuccess(resp, row)"
                        accept=".json">
                        <el-button circle plain type="warning" :icon="Upload">
                        </el-button>
                    </el-upload>
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

<style lang="less" scoped></style>