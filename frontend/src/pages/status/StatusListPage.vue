<script lang="ts" setup>
import { onMounted, ref } from 'vue'
import { queryRecordListService } from '@/api/userController'
import { formatTime } from '@/utils/format'

onMounted(() => {
    getRecordList()
})

const all_status = ['Accepted', 'Compile Error', 'Wrong Answer',
    'Memory Limit Exceeded', 'Time Limit Exceeded', 'Output Limit Exceeded', 'File Error', 'Nonzero Exit Status', 'Signalled', 'Internal Error']
const loading = ref(false)
const params = ref<API.QueryRecordListParams>({
    page: 1,
    page_size: 10,
    status: ''
})
const record_list = ref<API.Record[]>([])
const total = ref(0)
const getRecordList = async () => {
    loading.value = true
    const res = await queryRecordListService(params.value)
    console.log(res)
    record_list.value = res.data.data.data
    total.value = res.data.data.total
    loading.value = false
}
const handleCurrentChange = (page: number) => {
    params.value.page = page
    getRecordList()
}

</script>

<template>
    <page-container title="Status">
        <template #extra>
            <el-select v-model="params.status" clearable placeholder="status" style="width: 240px">
                <el-option v-for="item in all_status" :key="item" :label="item" :value="item" />
            </el-select>
        </template>

        <el-table v-loading="loading" :data="record_list" stripe="true">
            <el-table-column label="When" prop="created_at" width="200" align="center">
                <template #default="{ row }">
                    {{ formatTime(row.created_at) }}
                </template>
            </el-table-column>

            <el-table-column label="ID" prop="id" width="150" align="center">
                <template #default="{ row }">
                    <el-link type="primary" :underline="false" @click="
                        $router.push({
                            path: `/status/${row.id}`
                        })
                        ">{{ row.id }}</el-link>
                </template>
            </el-table-column>
            <el-table-column label="Status" prop="status" align="center">
                <template #default="{ row }">
                    <el-tag v-if="row.status === 'Accepted'" type="success">{{ row.status }}</el-tag>
                    <el-tag v-else-if="row.status === 'Compile Error'" type="warning">{{ row.status }}</el-tag>
                    <el-tag v-else type="danger">{{ row.status }}</el-tag>
                </template>
            </el-table-column>
            <el-table-column label="Problem" prop="problem_name" align="center">
                <template #default="{ row }">
                    {{ row.problem_name }}
                </template>
            </el-table-column>
            <el-table-column label="Language" prop="lang" align="center">
                <template #default="{ row }">
                    {{ row.lang }}
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
    </page-container>
</template>

<style lang="scss" scoped></style>