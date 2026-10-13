<script lang="ts" setup>
import { onMounted, ref } from 'vue'
import { Delete, Edit, Open, TurnOff } from '@element-plus/icons-vue'
import {
    queryNoticeListService,
    queryNoticeDetailService,
    addNoticeService,
    updateNoticeService,
    deleteNoticeService,
    publishNoticeService,
    hideNoticeService
} from '@/api/noticeController'
import { formatTime } from '@/utils/format'

onMounted(() => {
    getNoticeList()
})

const loading = ref(false)
const total = ref(0)
const noticeList = ref<API.Notice[]>([])
const params = ref<API.QueryNoticeListParams>({
    page: 1,
    page_size: 10,
})
const getNoticeList = async () => {
    loading.value = true
    const resp = await queryNoticeListService(params.value)
    total.value = Number(resp.data.data.total)
    noticeList.value = resp.data.data.list || []
    loading.value = false
}
const handleCurrentChange = (page: number) => {
    params.value.page = page
    getNoticeList()
}

// 创建 / 编辑抽屉
const drawerVisible = ref(false)
const formRef = ref()
const editMode = ref<'create' | 'update'>('create')
const form = ref<API.UpdateNoticeForm>({
    id: 0,
    title: '',
    content: '',
})
const rules = {
    title: [
        { required: true, message: '请输入公告标题', trigger: 'blur' },
    ],
    content: [
        { required: true, message: '请输入公告内容', trigger: 'blur' },
    ],
}
const onAddNotice = () => {
    editMode.value = 'create'
    form.value = { id: 0, title: '', content: '' }
    drawerVisible.value = true
}
const onEdit = async (row: API.Notice) => {
    editMode.value = 'update'
    const resp = await queryNoticeDetailService(Number(row.id))
    const notice = resp.data.data.notice
    form.value = {
        id: Number(notice.id),
        title: notice.title,
        content: notice.content,
    }
    drawerVisible.value = true
}
const onSubmit = async () => {
    await formRef.value.validate()
    let resp
    if (editMode.value === 'create') {
        resp = await addNoticeService({ title: form.value.title, content: form.value.content })
    } else {
        resp = await updateNoticeService(form.value)
    }
    if (resp.data.code === 0) {
        ElMessage.success(editMode.value === 'create' ? '创建成功' : '更新成功')
        drawerVisible.value = false
        getNoticeList()
    } else {
        ElMessage.error(resp.data.message || '提交失败')
    }
}

// 删除
const onDelete = async (row: API.Notice) => {
    ElMessageBox.confirm('确定要删除该公告吗？', '提示', {
        type: 'warning',
        confirmButtonText: '确定',
        cancelButtonText: '取消',
    }).then(async () => {
        const resp = await deleteNoticeService(Number(row.id))
        if (resp.data.code === 0) {
            ElMessage.success('删除成功')
            getNoticeList()
        } else {
            ElMessage.error(resp.data.message || '删除失败')
        }
    }).catch(() => { })
}

// 发布 / 隐藏
const onToggleStatus = async (row: API.Notice) => {
    const isPublished = row.status === 1
    const resp = isPublished
        ? await hideNoticeService(Number(row.id))
        : await publishNoticeService(Number(row.id))
    if (resp.data.code === 0) {
        ElMessage.success(isPublished ? '已隐藏' : '已发布')
        getNoticeList()
    } else {
        ElMessage.error(resp.data.message || '操作失败')
    }
}
</script>

<template>
    <page-container title="公告管理">
        <template #extra>
            <el-button type="primary" @click="onAddNotice">创建公告</el-button>
        </template>

        <!-- 表格区域 -->
        <el-table v-loading="loading" :data="noticeList">
            <el-table-column label="#" prop="id" width="80">
                <template #default="{ row }">
                    {{ row.id }}
                </template>
            </el-table-column>
            <el-table-column label="标题" prop="title" min-width="240">
                <template #default="{ row }">
                    {{ row.title }}
                </template>
            </el-table-column>
            <el-table-column label="创建时间" prop="createAt" width="180">
                <template #default="{ row }">
                    {{ formatTime(row.createAt) }}
                </template>
            </el-table-column>
            <el-table-column label="创建人" prop="createBy" width="120">
                <template #default="{ row }">
                    {{ row.createByName || row.createBy }}
                </template>
            </el-table-column>
            <el-table-column label="状态" width="100">
                <template #default="{ row }">
                    <el-tag v-if="row.status === 1" type="success">已发布</el-tag>
                    <el-tag v-else type="info">已隐藏</el-tag>
                </template>
            </el-table-column>
            <el-table-column label="操作" width="180">
                <template #default="{ row }">
                    <!-- 编辑 -->
                    <el-button circle plain type="primary" :icon="Edit" @click="onEdit(row)"></el-button>
                    <!-- 发布 / 隐藏 -->
                    <el-button v-if="row.status === 0" circle plain type="info" :icon="TurnOff"
                        @click="onToggleStatus(row)"></el-button>
                    <el-button v-else circle plain type="success" :icon="Open"
                        @click="onToggleStatus(row)"></el-button>
                    <!-- 删除 -->
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

        <!-- 创建 / 编辑公告抽屉 -->
        <el-drawer v-model="drawerVisible" size="520px" :with-header="false">
            <div class="drawer-header">
                <strong>{{ editMode === 'create' ? '创建公告' : '编辑公告' }}</strong>
                <span class="sub">{{ editMode === 'create' ? '发布一条新的平台公告' : '修改公告标题和内容' }}</span>
            </div>

            <el-form class="form" ref="formRef" :model="form" :rules="rules" label-position="top" size="large">
                <el-form-item label="标题" prop="title">
                    <el-input v-model="form.title" placeholder="请输入公告标题"></el-input>
                </el-form-item>
                <el-form-item label="内容" prop="content">
                    <el-input v-model="form.content" type="textarea" :rows="12" resize="none"
                        placeholder="请输入公告内容"></el-input>
                </el-form-item>
            </el-form>

            <div class="drawer-footer">
                <el-button size="large" @click="drawerVisible = false">取消</el-button>
                <el-button size="large" type="primary" @click="onSubmit">确认</el-button>
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
