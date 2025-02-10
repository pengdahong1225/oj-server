<script setup lang="ts">
import { ref } from 'vue'
import { updateProblemService, queryProblemDetailService, getProblemTagListService } from '@/api/problemController'
import type { Problem } from '@/types/problem.ts'

const level_list = [
    { label: '简单', value: 1 },
    { label: '中等', value: 2 },
    { label: '困难', value: 3 }
]
const drawerVisible = ref(false)
const formRef = ref()

const formModel = ref<Problem>({
    id: 0,
    title: '',
    description: '',
    level: 0,
    tags: [],
})
const emit = defineEmits(['success'])
const onSubmit = async () => {
    await formRef.value.validate()
    const isEdit = false
    if (isEdit) {
        const res = await updateProblemService(formModel.value)
        console.log(res)
    } else {
        const res = await updateProblemService(formModel.value)
        console.log(res)
    }
    drawerVisible.value = false
    emit('success', !isEdit)
}

const open = (row: Problem) => {
    drawerVisible.value = true
    formModel.value = row
    if (formModel.value.id > 0) {
        getProblemDetail(formModel.value.id)
    }
    getProblemTagList()
}
defineExpose({
    open
})

const getProblemDetail = async (id: number) => {
    const res = await queryProblemDetailService(id)
    console.log(res)
    formModel.value.id = res.data.data.id
    formModel.value.title = res.data.data.title
    formModel.value.description = res.data.data.description
    formModel.value.level = res.data.data.level
    formModel.value.tags = res.data.data.tags
    formModel.value.status = res.data.data.status
    formModel.value.config = JSON.stringify(res.data.data.config)
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
            <strong>{{ formModel.id ? '编辑题目' : '添加题目' }}</strong>
        </template>

        <template #default>
            <el-form class="form" ref="formRef" :model="formModel" label-width="150px" size="large">
                <el-form-item label="标题" prop="title">
                    <el-input v-model="formModel.title" placeholder="请输入题目标题" style="width: 50%"></el-input>
                </el-form-item>
                <el-form-item label="Level" prop="level">
                    <el-select size="large" v-model="formModel.level" placeholder="Select" style="width: 240px">
                        <el-option v-for="item in level_list" :key="item.value" :label="item.label"
                            :value="item.value" />
                    </el-select>
                </el-form-item>
                <el-form-item label="标签" prop="tags">
                    <el-select v-model="formModel.tags" multiple filterable allow-create default-first-option
                        placeholder="Select" style="width: 240px">
                        <el-option v-for="item in tag_list" :key="item" :label="item" :value="item" />
                    </el-select>
                </el-form-item>
                <el-form-item label="题目详情" prop="description">
                    <el-input type="textarea" :rows="10" resize="none" v-model="formModel.description"
                        placeholder="请描述题目详情" style="width: 100%"></el-input>
                </el-form-item>
                <el-form-item label="题目配置" prop="config">
                    <el-input type="textarea" :rows="10" resize="none" v-model="formModel.config" placeholder="请提供题目配置"
                        style="width: 100%"></el-input>
                </el-form-item>
                <el-form-item style="">
                    <el-button @click="drawerVisible = false">取消</el-button>
                    <el-button type="primary" @click="onSubmit"> 确认 </el-button>
                </el-form-item>
            </el-form>
        </template>
    </el-drawer>
</template>

<style lang="scss" scoped>
.form {
    width: 100%;

    // 表单中的元素居中，占78%宽度
    .el-form-item {
        width: 80%;
        margin: 10px 0px 10px 0px;
    }
}
</style>