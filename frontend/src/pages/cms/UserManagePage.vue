<script lang="ts" setup>
import { onMounted, ref } from 'vue'
import { Delete, Edit, Search } from '@element-plus/icons-vue'
import { queryUserListService, createUserService } from '@/api/userController'
import { formatTime } from '@/utils/format'

onMounted(() => {
    getUserList()
})

const loading = ref(false)
const total = ref(0)
const userList = ref<API.UserInfo[]>([])
const params = ref(<API.QueryUserListParams>{
    page: 1,
    page_size: 10, // page_size默认为10
    keyword: '',
})
const getUserList = async () => {
    loading.value = true
    const resp = await queryUserListService(params.value)
    total.value = Number(resp.data.data.total)
    userList.value = resp.data.data.data
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

// 创建用户
const drawerVisible = ref(false)
const formRef = ref()
const createForm = ref<API.CreateUserForm>({
    mobile: '',
    password: '',
    username: '',
    email: '',
    role: 0,
})
const roleList = [
    { label: '用户', value: 0 },
    { label: '管理员', value: 1 },
]
const rules = {
    mobile: [
        { required: true, message: '请输入手机号', trigger: 'blur' },
        { pattern: /^1[3-9]\d{9}$/, message: '手机号格式错误', trigger: 'blur' },
    ],
    password: [
        { required: true, message: '请输入密码', trigger: 'blur' },
    ],
    username: [
        { required: true, message: '请输入用户名', trigger: 'blur' },
    ],
}
const onAddUser = () => {
    createForm.value = {
        mobile: '',
        password: '',
        username: '',
        email: '',
        role: 0,
    }
    drawerVisible.value = true
}
const onSubmit = async () => {
    await formRef.value.validate()
    const resp = await createUserService(createForm.value)
    if (resp.data.code === 0) {
        ElMessage.success('创建成功')
        drawerVisible.value = false
        getUserList()
    } else {
        ElMessage.error(resp.data.message || '创建失败')
    }
}

</script>

<template>
    <page-container title="用户管理">
        <template #extra>
            <el-button type="primary" @click="onAddUser">创建用户</el-button>
        </template>

        <!-- 表单区域 inline属性代表表单元素都在一行显示 -->
        <el-form inline>
            <el-form-item style="margin-right: 15px;">
                <el-input v-model="params.keyword" style="width: 240px" placeholder="用户名"
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
            <el-table-column label="用户名" prop="title">
                <template #default="{ row }">
                    {{ row.username }}
                </template>
            </el-table-column>
            <el-table-column label="手机号" prop="mobile">
                <template #default="{ row }">
                    {{ row.mobile }}
                </template>
            </el-table-column>
            <el-table-column label="邮箱">
                <template #default="{ row }">
                    {{ row.email || 'nil' }}
                </template>
            </el-table-column>
            <el-table-column label="注册时间">
                <template #default="{ row }">
                    {{ formatTime(row.createAt) }}
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

        <!-- 创建用户抽屉 -->
        <el-drawer v-model="drawerVisible" size="420px" :with-header="false">
            <div class="drawer-header">
                <strong>创建用户</strong>
                <span class="sub">创建一个新的平台账号</span>
            </div>

            <el-form class="form" ref="formRef" :model="createForm" :rules="rules" label-position="top" size="large">
                <el-form-item label="用户名" prop="username">
                    <el-input v-model="createForm.username" placeholder="请输入用户名"></el-input>
                </el-form-item>
                <el-form-item label="手机号" prop="mobile">
                    <el-input v-model="createForm.mobile" placeholder="请输入手机号"></el-input>
                </el-form-item>
                <el-form-item label="密码" prop="password">
                    <el-input v-model="createForm.password" type="password" show-password placeholder="请输入密码"></el-input>
                </el-form-item>
                <el-form-item label="邮箱" prop="email">
                    <el-input v-model="createForm.email" placeholder="请输入邮箱（选填）"></el-input>
                </el-form-item>
                <el-form-item label="角色" prop="role">
                    <el-select v-model="createForm.role" placeholder="请选择角色" style="width: 100%">
                        <el-option v-for="item in roleList" :key="item.value" :label="item.label" :value="item.value" />
                    </el-select>
                </el-form-item>
            </el-form>

            <div class="drawer-footer">
                <el-button size="large" @click="drawerVisible = false">取消</el-button>
                <el-button size="large" type="primary" @click="onSubmit">创建</el-button>
            </div>
        </el-drawer>
    </page-container>
</template>

<style lang="less" scoped>
.drawer-header {
    display: flex;
    flex-direction: column;
    gap: 4px;
    padding-bottom: 16px;
    border-bottom: 1px solid #f0f0f0;

    strong {
        font-size: 20px;
        color: #303133;
    }

    .sub {
        font-size: 13px;
        color: #909399;
    }
}

.form {
    padding-top: 24px;

    :deep(.el-form-item__label) {
        font-weight: 500;
        color: #303133;
    }
}

.drawer-footer {
    display: flex;
    justify-content: flex-end;
    gap: 12px;
    margin-top: 24px;
    padding-top: 16px;
    border-top: 1px solid #f0f0f0;

    .el-button {
        min-width: 96px;
    }
}
</style>