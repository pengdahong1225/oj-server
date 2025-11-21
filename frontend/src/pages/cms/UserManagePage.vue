<script lang="ts" setup>
import { onMounted, ref } from 'vue'
import { Delete, Edit, Search } from '@element-plus/icons-vue'
import { queryUserListService } from '@/api/userController'
import { formatTime } from '@/utils/format'

onMounted(() => {
    getUserList()
})

const loading = ref(false)
// 总数
const total = ref(0)
// 用户列表数据
const userList = ref<API.Problem[]>([])
// 分页请求参数
const params = ref(<API.QueryUserListParams>{
    page: 1,
    page_size: 10, // page_size默认为10
    keyword: '',
})
const getUserList = async () => {
    loading.value = true
    const resp = await queryUserListService(params.value)
    total.value = resp.data.data.total
    userList.value = resp.data.data.list
    loading.value = false
}
const handleCurrentChange = (page: number) => {
    params.value.page = page
    getUserList()
}
const onSearch = () => {
    params.value.page = 1
    getUserList()
}
const onReset = () => {
    params.value.page = 1
    params.value.keyword = ''
    getUserList()
}

const onEdit = (row: API.UserInfo) => {
    console.log(row)
}
const onDelete = async (id: number) => {
   console.log(id)
}

</script>

<template>
    <page-container title="用户管理">
        <!-- 表单区域 inline属性代表表单元素都在一行显示 -->
        <el-form inline>
            <el-form-item style="margin-right: 15px;">
                <el-input v-model="params.keyword" style="width: 240px" placeholder="nick name"
                    :suffix-icon="Search" />
            </el-form-item>
            <el-form-item>
                <el-button @click="onSearch" type="primary">搜索</el-button>
                <el-button @click="onReset">重置</el-button>
            </el-form-item>
        </el-form>

        <!-- 表格区域 -->
        <el-table v-loading="loading" :data="userList">
            <el-table-column label="#" prop="id" width="100">
                <template #default="{ row }">
                    {{ row.uid }}
                </template>
            </el-table-column>
            <el-table-column label="手机号码" prop="mobile">
                <template #default="{ row }">
                    {{ row.mobile }}
                </template>
            </el-table-column>
            <el-table-column label="用户名" prop="title">
                <template #default="{ row }">
                    {{ row.nickname }}
                </template>
            </el-table-column>
            <el-table-column label="邮箱">
                <template #default="{ row }">
                    {{ row.email || 'nil' }}
                </template>
            </el-table-column>
            <el-table-column label="注册时间">
                <template #default="{ row }">
                    {{ formatTime(row.create_at) }}
                </template>
            </el-table-column>
            <el-table-column label="角色">
                <template #default="{ row }">
                    <el-tag v-if="row.role === 1" type="danger">管理员</el-tag>
                    <el-tag v-else type="primary">用户</el-tag>
                </template>
            </el-table-column>
             <el-table-column label="操作">
                <template #default="{ row }">
                    <!-- 编辑 -->
                    <el-button circle plain type="primary" :icon="Edit" @click="onEdit(row)"></el-button>
                    <!-- 删除 -->
                    <el-button circle plain type="danger" :icon="Delete" @click="onDelete(row.uid)"></el-button>
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

<style lang="less" scoped></style>