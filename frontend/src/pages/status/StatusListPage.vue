<script lang="ts" setup>
import { onMounted, ref, computed } from 'vue'
import { queryRecordListService } from '@/api/userController'
import { formatTime } from '@/utils/format'

onMounted(() => {
    getRecordList()
})

const loading = ref(false)
const params = ref<API.QueryRecordListParams>({
    page: 1,
    page_size: 10
})
const record_list = ref<API.Record[]>([])
const total = ref(0)
const getRecordList = async () => {
    loading.value = true
    const resp = await queryRecordListService(params.value)
    total.value = resp.data.data.total
    record_list.value = resp.data.data.list
    loading.value = false
}
const handleCurrentChange = (page: number) => {
    params.value.page = page
    getRecordList()
}
const langMap: Record<string, string> = {
    c: 'C',
    cpp: 'C++',
    java: 'Java',
    python: 'Python',
    golang: 'Golang',
}
const record_list_with_computed = computed(() =>
   (record_list.value || []).map(row => ({
    ...row,
    lang: langMap[row?.lang ?? ''] ?? row?.lang ?? ''
  }))
)

</script>

<template>
    <div class="status-list-container">
        <el-table v-loading="loading" :data="record_list_with_computed" stripe>
            <el-table-column label="ID" prop="id" width="50" align="center">
                <template #default="{ row }">
                    <el-link type="primary" @click="
                        $router.push({
                            path: `/status/${row.id}`
                        })
                        ">{{ row.id }}</el-link>
                </template>
            </el-table-column>

            <el-table-column label="状态" prop="status" align="center">
                <template #default="{ row }">
                    <el-tag v-if="row.message === 'Accepted'" type="success">{{ row.message }}</el-tag>
                    <el-tag v-else-if="row.message === 'Compile Error'" type="warning">{{ row.message }}</el-tag>
                    <el-tag v-else type="danger">{{ row.message }}</el-tag>
                </template>
            </el-table-column>

            <el-table-column label="题目" prop="problem_name" align="center">
                <template #default="{ row }">
                    {{ (row.problem_name == null || row.problem_name == '') ? row.problem_id : row.problem_name }}
                </template>
            </el-table-column>

            <el-table-column label="语言" prop="lang" align="center">
                <template #default="{ row }">
                    {{ row.lang }}
                </template>
            </el-table-column>
            <el-table-column label="时间" prop="clock" align="center">
                <template #default="{ row }">
                    {{ row.clock }}
                </template>
            </el-table-column>
            <el-table-column label="内存" prop="memory" align="center">
                <template #default="{ row }">
                    {{ row.memory }}
                </template>
            </el-table-column>

            <el-table-column label="提交时间" prop="created_at" width="200" align="center">
                <template #default="{ row }">
                    {{ formatTime(row.created_at) }}
                </template>
            </el-table-column>
            <el-table-column label="提交者" prop="user_name" width="200" align="center">
                <template #default="{ row }">
                    {{ (row.user_name == null || row.user_name == '') ? row.uid : row.user_name }}
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
    </div>
</template>

<style lang="less" scoped>
// 增加左右两边边距
.status-list-container {
    padding-left: 50px;
    padding-right: 50px;
}

// 深度选择器，将链接颜色改为黑色
:deep(.el-link--primary) {
    color: #000 !important;
    
    &:hover {
        color: #333 !important;
    }
}
</style>