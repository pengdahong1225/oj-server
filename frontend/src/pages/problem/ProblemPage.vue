<script lang="ts" setup>
import { onMounted, ref, computed } from 'vue'
import { useRoute } from 'vue-router'
import { queryProblemDetailService, submitProblemService, queryResultService } from '@/api/problemController'
import { getCommentListService } from '@/api/commentController'
import { UploadFilled } from '@element-plus/icons-vue'
import { VAceEditor } from 'vue3-ace-editor'
import 'ace-builds/src-noconflict/mode-c_cpp'
import 'ace-builds/src-noconflict/mode-golang'
import 'ace-builds/src-noconflict/mode-python'
import 'ace-builds/src-noconflict/mode-java'
import 'ace-builds/src-noconflict/theme-github'
import 'ace-builds/src-noconflict/theme-xcode'
import { marked } from 'marked'

// 题目详细信息
onMounted(() => {
    // 获取题目信息
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

    // 获取评论列表
    getCommentList()
}
const btitle = computed(() => {
    return problem.value.id + '. ' + problem.value.title
})
const bdescription = computed(() => {
    return marked.parse(problem.value.description)
})

// 题目提交表单
const loading = ref(false)
const form = ref<SubmitForm>({
    code: '',
    lang: 'c',  // 默认使用c
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
    switch (res.data.message) {
        case 'running...':
            break
        case 'OK':
            clearInterval(timer)
            running.value = false
            result.value = res.data.data
            break
        default:
            clearInterval(timer)
            running.value = false
            console.log(res.data.message)
            break
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
        return 'You have solved this problem'
    } else {
        return 'failed'
    }
})
const showed = ref(false)
const showResult = () => {
    showed.value = true
    // 跳转
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
    highlightActiveLine: false,
    enableSnippets: true,
    wrap: true, // 自动换行
    showGutter: true, // 显示行号
});

// 评论区数据
const comment_query_params = <API.QueryCommentListParams>{}
const comment_text = ref('')
const root_comment_list = ref<API.Comment[]>([])
const getCommentList = async () => {
    comment_query_params.obj_id = problem.value.id
    comment_query_params.cursor = 0

    const res = await getCommentListService(comment_query_params)
    console.log(res)
}

</script>

<template>
    <el-card class="container" v-loading="running">
        <!-- 标题区域 -->
        <template #header>
            <div class="header">
                <div class="header-title">{{ btitle }}</div>
                <div class="header-foot">
                    <div class="header-foot-level">
                        <el-tag v-if="problem.level === 1" type="primary">简单</el-tag>
                        <el-tag v-else-if="problem.level === 2" type="warning">中等</el-tag>
                        <el-tag v-else type="danger">困难</el-tag>
                    </div>
                    <div class="header-foot-tags">
                        {{ problem.tags }}
                    </div>
                </div>
            </div>
        </template>

        <!-- 题目详情区域 -->
        <div class="description-box">
            <div class="description" v-html="bdescription"></div>
            <div class="statistic">
                <div>通过次数 6.1M</div>
                <el-divider direction="vertical" border-style="dashed" />
                <div>提交次数 11.1M</div>
                <el-divider direction="vertical" border-style="dashed" />
                <div>通过率 54.6%</div>
            </div>
        </div>

        <el-divider />

        <!-- 题目编辑与提交区域 -->
        <div class="edit_submit">
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

            <VAceEditor class="editor" v-model:value="form.code" :lang="lang_computed" :theme="theme"
                :options="editorOptions" />

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
        </div>

        <el-divider />

        <!-- 评论区 -->
        <div class="comment-area">
            <!-- 输入评论 -->
            <div class="comment-input">
                <el-input v-model="comment_text" type="textarea" show-word-limit="true" placeholder="请输入评论..."
                    resize="none" :autosize="{ minRows: 5 }" maxlength="1000"
                    input-style="font-size: 18px; border: none; outline: none; overflow: hidden;" />
                <div style="display: flex; justify-content: flex-end;">
                    <el-button type="success" :disabled="!comment_text" auto-insert-space
                        style="margin-top: 5px;">评论</el-button>
                </div>
            </div>
            <el-card shadow="never">
                <!-- 评论区规则 -->
                <template #header>
                    <div class="card-header">
                        <strong>💡评论区规则</strong>
                        <p>1. 请不要在评论区发表题解！</p>
                        <p>2. 评论区可以发表关于对翻译的建议、对题目的疑问及其延伸讨论。</p>
                        <p>3. 如果你需要整理题解思路，获得反馈从而进阶提升，可以去题解区进行。</p>
                    </div>
                </template>
                <!-- 评论列表 -->
                <div class="comment-list">
                    <!-- 顶层评论 -->
                    <!-- <RootCommentItem></RootCommentItem> -->
                </div>
            </el-card>
        </div>
    </el-card>
</template>

<style lang="less" scoped>
.container {
    width: 90%;
    height: auto;
    margin: auto;

    .header {
        display: block;

        .header-title {
            font-family: "Arial", "Helvetica", sans-serif;
            font-size: 25px;
            font-weight: bold;
            line-height: 1.6; // 行间距
        }

        .header-foot {
            display: flex;

            .header-foot-level {
                margin-right: 5px;
            }

            .header-foot-tags {
                margin-left: 5px;
            }
        }
    }

    .description-box {
        display: block;
        padding-left: 20px;

        .description {
            align-items: flex-start;
            margin-bottom: 20px;
            font-size: 17px;
        }

        .statistic {
            display: flex;
        }
    }

    .edit_submit {
        margin-top: 20px;
        width: 100%;

        .editor {
            height: 350px;
            width: 100%;
        }

        .submit-area-foot {
            margin: auto;
            margin-top: 10px;
            width: 100%;
            display: flex;

            .result-area {
                width: 50%;
            }

            .submit-btn-area {
                width: 50%;
                display: flex;
                justify-content: flex-end;
            }
        }
    }

    .comment-area {
        margin-top: 20px;
        width: 100%;
        height: auto;
        display: block;

        .comment-input {
            height: auto;
            margin-bottom: 5px;
        }

        .comment-list {
            margin-left: 0px;
            margin-right: 0px;
        }
    }
}
</style>