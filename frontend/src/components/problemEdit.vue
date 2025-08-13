<script setup lang="ts">
import { ref } from 'vue'
import { addProblemService, updateProblemService, queryProblemDetailService, getProblemTagListService } from '@/api/problemController'

// 父子组件通信数据
const open = (mode: string, row: API.Problem) => {
    drawerVisible.value = true
    edit_mode.value = mode
    if (mode === 'update') {
        update_form.value = row
        getProblemDetail(update_form.value.id)
    }
    getProblemTagList()
}
defineExpose({
    open
})
const emit = defineEmits(['success'])
const onSubmit = async () => {
    await formRef.value.validate()

    let resp
    if (edit_mode.value === 'create') {
        resp = await addProblemService(create_form.value)
    } else {
        resp = await updateProblemService(update_form.value)
    }
    console.log(resp)
    if (resp.data.data.code === 0) {
        ElMessage.success('提交成功')
    } else {
        ElMessage.error(resp.data?.data?.message || '提交失败')
    }

    drawerVisible.value = false
    emit('success', edit_mode.value)
}

const level_list = [
    { label: '简单', value: 1 },
    { label: '中等', value: 2 },
    { label: '困难', value: 3 }
]
const drawerVisible = ref(false)
const formRef = ref()

// 默认创建
let edit_mode = ref('create')
const create_form = ref<API.CreateProblemForm>({
    title: '',
    level: 0,
    tags: [],
    description: '',
})
const update_form = ref<API.UpdateProblemForm>({
    id:  0,
    title: '',
    level: 0,
    tags: [],
    description: '',
})

const getProblemDetail = async (id: number) => {
    const resp = await queryProblemDetailService(id)
    console.log(resp)
    update_form.value.id = resp.data.data.id
    update_form.value.title = resp.data.data.title
    update_form.value.description = resp.data.data.description
    update_form.value.level = resp.data.data.level
    update_form.value.tags = resp.data.data.tags
}

const tag_list = ref<string[]>([])
const getProblemTagList = async () => {
    const res = await getProblemTagListService()
    console.log(res)
    if (Array.isArray(res.data.data) && res.data.data.length > 0) {
        tag_list.value = res.data.data
    }
}
</script>

<template>
    <el-drawer v-model="drawerVisible" size="45%">
        <template #header>
            <strong v-if="edit_mode == 'create'"> 新增题目 </strong>
            <strong v-else> 编辑题目 </strong>
        </template>

        <template #default>
            <el-form v-if="edit_mode == 'create'" class="form" ref="formRef" :model="update_form" label-width="150px" size="large">
                <el-form-item label="标题" prop="title">
                    <el-input v-model="update_form.title" placeholder="请输入题目标题" style="width: 50%"></el-input>
                </el-form-item>
                <el-form-item label="Level" prop="level">
                    <el-select size="large" v-model="update_form.level" placeholder="Select" style="width: 240px">
                        <el-option v-for="item in level_list" :key="item.value" :label="item.label"
                            :value="item.value" />
                    </el-select>
                </el-form-item>
                <el-form-item label="标签" prop="tags">
                    <el-select v-model="update_form.tags" multiple filterable allow-create default-first-option
                        placeholder="Select" style="width: 240px">
                        <el-option v-for="item in tag_list" :key="item" :label="item" :value="item" />
                    </el-select>
                </el-form-item>
                <el-form-item label="题目详情" prop="description">
                    <el-input type="textarea" :rows="10" resize="none" v-model="update_form.description"
                        placeholder="请描述题目详情" style="width: 100%"></el-input>
                </el-form-item>
                <el-form-item style="">
                    <el-button @click="drawerVisible = false">取消</el-button>
                    <el-button type="primary" @click="onSubmit"> 确认 </el-button>
                </el-form-item>
            </el-form>
            <el-form v-else class="form" ref="formRef" :model="update_form" label-width="150px" size="large">
                <el-form-item label="标题" prop="title">
                    <el-input v-model="update_form.title" placeholder="请输入题目标题" style="width: 50%"></el-input>
                </el-form-item>
                <el-form-item label="Level" prop="level">
                    <el-select size="large" v-model="update_form.level" placeholder="Select" style="width: 240px">
                        <el-option v-for="item in level_list" :key="item.value" :label="item.label"
                            :value="item.value" />
                    </el-select>
                </el-form-item>
                <el-form-item label="标签" prop="tags">
                    <el-select v-model="update_form.tags" multiple filterable allow-create default-first-option
                        placeholder="Select" style="width: 240px">
                        <el-option v-for="item in tag_list" :key="item" :label="item" :value="item" />
                    </el-select>
                </el-form-item>
                <el-form-item label="题目详情" prop="description">
                    <el-input type="textarea" :rows="10" resize="none" v-model="update_form.description"
                        placeholder="请描述题目详情" style="width: 100%"></el-input>
                </el-form-item>
                <el-form-item style="">
                    <el-button @click="drawerVisible = false">取消</el-button>
                    <el-button type="primary" @click="onSubmit"> 确认 </el-button>
                </el-form-item>
            </el-form>

        </template>
    </el-drawer>
</template>

<style lang="less" scoped>
.form {
    width: 100%;

    // 表单中的元素居中，占78%宽度
    .el-form-item {
        width: 80%;
        margin: 10px 0px 10px 0px;
    }
}
</style>