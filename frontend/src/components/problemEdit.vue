<script setup lang="ts">
import { ref } from 'vue'
import { updateProblemService } from '@/api/problem.ts'
import type { Problem } from '@/types/problem.ts'

const level_list = [
    { label: '简单', value: 1 },
    { label: '中等', value: 2 },
    { label: '困难', value: 3 }
]
const drawerVisible = ref(false)
const formRef = ref()
const props = defineProps({
    formModel: {
        type: Object as () => Problem,
        default: () => ({
            title: '',
            description: '',
            tags: [],
            level: 0
        })
    }
})
const emit = defineEmits(['success'])
const rules = {}
const onSubmit = async () => {
    await formRef.value.validate()
    const isEdit = false
    if (isEdit) {
        // await artEditChannelService(formModel.value)
        // ElMessage.success('编辑成功')
    } else {
        // await artAddChannelService(formModel.value)
        // ElMessage.success('添加成功')
    }
    drawerVisible.value = false
    emit('success', !isEdit)
}

// 组件对外暴露一个方法 open，基于open传来的参数，区分添加还是编辑
// open({})  => 表单无需渲染，说明是添加
// open({ id, cate_name, ... })  => 表单需要渲染，说明是编辑
// open调用后，可以打开弹窗
const open = (row: Problem) => {
    drawerVisible.value = true
    //   formModel.value = row
}

// 向外暴露方法
defineExpose({
    open
})
</script>

<template>
    <el-drawer v-model="drawerVisible" size="45%">
        <template #header>
            <strong>{{ formModel?.id ? '编辑题目' : '添加题目' }}</strong>
        </template>

        <template #default>
            <el-form class="form" ref="formRef" :model="formModel" :rules="rules" label-width="150px" size="large">
                <el-form-item label="标题" prop="title">
                    <el-input v-model="formModel.title" placeholder="请输入题目标题" style="width: 50%"></el-input>
                </el-form-item>
                <el-form-item label="Level" prop="level">
                    <el-select size="large" v-model="formModel.level" placeholder="Select" style="width: 240px">
                        <el-option v-for="item in level_list" :key="item.value" :label="item.label"
                            :value="item.value" />
                    </el-select>
                </el-form-item>
                <el-form-item label="题目详情" prop="description">
                    <el-input type="textarea" :rows="10" resize="none" v-model="formModel.description"
                        placeholder="请描述题目详情" style="width: 100%"></el-input>
                </el-form-item>
                <el-form-item label="Sample Input" prop="sample_input">
                    <el-input v-model="formModel.description" placeholder="输入样例"></el-input>
                </el-form-item>
                <el-form-item label="Sample Output" prop="sample_output">
                    <el-input v-model="formModel.description" placeholder="输出样例"></el-input>
                </el-form-item>
                <el-form-item label="Hint" prop="hint">
                    <el-input type="textarea" :rows="10" v-model="formModel.description" placeholder="代码示例"
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