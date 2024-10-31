<script lang="ts" setup>
import { onMounted, ref, computed } from 'vue'
import { useRoute } from 'vue-router'
import { queryProblemDetailService, submitProblemService } from '@/api/problem'
import type { Problem, SubmitForm } from '@/types/problem'
import { UploadFilled } from '@element-plus/icons-vue'
import { VAceEditor } from 'vue3-ace-editor'
import 'ace-builds/src-noconflict/mode-c_cpp'
import 'ace-builds/src-noconflict/mode-golang'
import 'ace-builds/src-noconflict/mode-python'
import 'ace-builds/src-noconflict/mode-java'
import 'ace-builds/src-noconflict/theme-github'
import 'ace-builds/src-noconflict/theme-xcode'

// 题目详细信息
onMounted(() => {
    getProblemDetail()
})
const route = useRoute()
const problem = ref<Problem>({
    id: 0,
    title: '',
    description: '',
    level: 0,
    tags: [],
    status: 1
})
const getProblemDetail = async () => {
    const res = await queryProblemDetailService(Number(route.params.id))
    problem.value.id = res.data.data.id
    problem.value.title = res.data.data.title
    problem.value.description = res.data.data.description
    problem.value.level = res.data.data.level
    problem.value.tags = res.data.data.tags
    problem.value.status = res.data.data.status
}

// 题目提交表单
const loading = ref(false)
const form = ref<SubmitForm>({
    code: '',
    lang: 'C',  // 默认使用C
    problem_id: 0,
    title: ''
})
const onReset = () => {
    form.value.code = ''
}
const onSubmit = async () => {
    loading.value = true
    form.value.problem_id = problem.value.id
    form.value.title = problem.value.title
    const res = await submitProblemService(form.value)
    console.log(res)
    loading.value = false
}

// ACE主题和配置
const lang_list = [
    { label: 'C', value: 'c' },
    { label: 'C++', value: 'c_cpp' },
    { label: 'Java', value: 'java' },
    { label: 'Python', value: 'python' },
    { label: 'Go', value: 'golang' },
]
const lang_computed = computed(() => {
    if (form.value.lang === 'c') {
        return 'c_cpp'
    } else {
        return form.value.lang
    }
})
const theme = ref('xcode');
const editorOptions = ref({
    fontSize: '20px',
    showPrintMargin: false,
    enableBasicAutocompletion: true,
    enableLiveAutocompletion: true,
    highlightActiveLine: true,
    enableSnippets: true,
});

</script>

<template>
    <div class="container">
        <!-- 左边题目描述和编辑区域 -->
        <div class="left">
            <el-card class="problem-description" shadow="hover">
                <template #header>
                    <strong>{{ problem.title }}</strong>
                </template>

                <div style="background: rgb(251, 251, 251); padding: 15px;">
                    <descriptionItem class="descriptionItem" title="Description" :content="problem.description"
                        style="width: 100%;">
                    </descriptionItem>
                    <descriptionItem class="descriptionItem" title="Input" :content="problem.description"
                        style="width: 100%;">
                    </descriptionItem>
                    <descriptionItem class="descriptionItem" title="Output" :content="problem.description"
                        style="width: 100%;">
                    </descriptionItem>
                    <descriptionItem class="descriptionItem" title="Sample Input" :content="problem.description"
                        style="width: 100%;">
                    </descriptionItem>
                    <descriptionItem class="descriptionItem" title="Sample Output" :content="problem.description"
                        style="width: 100%;">
                    </descriptionItem>
                    <descriptionItem class="descriptionItem" title="Hint" :content="problem.description"
                        style="width: 100%;">
                    </descriptionItem>
                </div>
            </el-card>

            <!-- 题目提交区域 -->
            <el-card class="submit-area" shadow="hover">
                <template #header>
                    <el-form inline class="form">
                        <el-form-item style="margin-right: 10px;">
                            <el-select size="large" v-model="form.lang" placeholder="Select" style="width: 240px">
                                <el-option v-for="item in lang_list" :key="item.label" :label="item.label"
                                    :value="item.value" />
                            </el-select>
                        </el-form-item>
                        <el-form-item>
                            <el-button size="large" @click="onReset">重置</el-button>
                        </el-form-item>
                    </el-form>
                </template>

                <VAceEditor v-model:value="form.code" :lang="lang_computed" :theme="theme" :options="editorOptions"
                    style="height: 500px; width: 100%;" />

                <el-button class="submit-btn" type="warning" :loading="loading" @click="onSubmit">
                    <el-icon><UploadFilled /></el-icon>Submit
                </el-button>
            </el-card>

        </div>

        <!-- 右边题目information区域 -->
        <div class="right">
            <el-card shadow="hover">
                <template #header>
                    <strong>Information</strong>
                </template>
                <el-descriptions :column="1" border>
                    <el-descriptions-item label="ID">1</el-descriptions-item>
                    <el-descriptions-item label="Time Limit">18100000000</el-descriptions-item>
                    <el-descriptions-item label="Memory Limit">Suzhou</el-descriptions-item>
                    <el-descriptions-item label="IO Mode">Suzhou</el-descriptions-item>
                    <el-descriptions-item label="Created By">Suzhou</el-descriptions-item>
                    <el-descriptions-item label="Level">Suzhou</el-descriptions-item>
                    <el-descriptions-item label="Tags">
                        <el-tag size="small">School</el-tag>
                    </el-descriptions-item>
                </el-descriptions>
            </el-card>
        </div>
    </div>
</template>

<style lang="scss" scoped>
.container {
    display: flex;
    width: 80%;
    margin: auto;
    align-items: flex-start;

    /* 确保卡片顶部对齐，而不会拉伸高度 */
    .left {
        width: 80%;
        margin-right: 5px;

        .descriptionItem {
            margin-bottom: 20px;
        }

        .submit-area {
            margin-top: 20px;
            width: 100%;
            .submit-btn{
                margin: 5px 0 5px 5px;
                // 靠右
                float: right;
            }
        }
    }

    .right {
        width: 20%;
        margin-left: 5px;
    }
}
</style>