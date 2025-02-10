<script lang="ts" setup>
import { onMounted, ref, computed } from 'vue'
import { useRoute } from 'vue-router'
import { queryProblemDetailService, submitProblemService, queryResultService } from '@/api/problemController'
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
    console.log(res)
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

    // 轮训判题结果
    running.value = true
    timer = setInterval(() => {
        checkResult()
    }, 500)
}

// 轮训判题结果
let timer = 0
const running = ref(false)
const checkResult = async () => {
    const res = await queryResultService(Number(route.params.id))
    console.log(res)
    if (res.data.data.message !== 'OK') {
        clearInterval(timer)
        running.value = false

        result.value = res.data.data
        console.log(result.value)
    }
}
const result = ref('')
const btype = computed(() => {
    if (result.value === 'Accepted') {
        return 'success'
    } else {
        return 'danger'
    }
})
const tmsg = computed(() => {
    if (result.value === 'Accepted') {
        return 'You have solved the problem'
    } else {
        return 'failed'
    }
})
const showed = ref(false)
const showResult = () => {
    showed.value = true
    // 跳转status
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
                    <descriptionItem class="descriptionItem" title="Hint" :content="problem.description"
                        style="width: 100%;">
                    </descriptionItem>
                </div>
            </el-card>

            <!-- 题目提交区域 -->
            <el-card class="submit-area" shadow="hover" v-loading="running">
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
                    style="height: 350px; width: 100%;" />

                <div class="submit-area-foot">
                    <div class="result-area">
                        <div v-if="result !== ''">
                            <el-button v-if="!showed" class="show-result-bnt" :type="btype" @click="showResult">
                                {{ result }}
                            </el-button>
                            <el-tag v-else :type="btype" size="large">
                                <template #default>
                                    <strong>{{ tmsg }}</strong>
                                </template>
                            </el-tag>
                        </div>
                    </div>

                    <div class="submit-btn-area">
                        <el-button class="submit-btn" type="warning" :loading="loading" @click="onSubmit">
                            <el-icon>
                                <UploadFilled />
                            </el-icon>Submit
                        </el-button>
                    </div>
                </div>
            </el-card>

        </div>

        <!-- 右边题目information区域 -->
        <div class="right">
            <el-card shadow="hover">
                <template #header>
                    <strong>Information</strong>
                </template>
                <el-descriptions :column="1" border>
                    <el-descriptions-item label="ID">{{ problem.id }}</el-descriptions-item>
                    <el-descriptions-item label="Time Limit">18100000000</el-descriptions-item>
                    <el-descriptions-item label="Memory Limit">Suzhou</el-descriptions-item>
                    <el-descriptions-item label="IO Mode">IO</el-descriptions-item>
                    <el-descriptions-item label="Created By">{{ problem.create_by }}</el-descriptions-item>
                    <el-descriptions-item label="Level">
                        <el-tag v-if="problem.level === 1" type="primary">简单</el-tag>
                        <el-tag v-else-if="problem.level === 2" type="warning">中等</el-tag>
                        <el-tag v-else type="danger">困难</el-tag>
                    </el-descriptions-item>
                    <el-descriptions-item label="Tags">
                        <el-tag size="small">{{ problem.tags }}</el-tag>
                    </el-descriptions-item>
                </el-descriptions>
            </el-card>
        </div>
    </div>
</template>

<style lang="less" scoped>
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

            .submit-area-foot {
                margin: auto;
                margin-top: 10px;
                width: 95%;
                display: flex;

                .result-area {
                    width: 50%;
                }

                .submit-btn-area {
                    width: 50%;

                    .submit-btn {
                        float: right;
                    }
                }
            }
        }
    }

    .right {
        width: 20%;
        margin-left: 5px;
    }
}
</style>