<script lang="ts" setup>
import { onMounted, ref, computed } from 'vue'
import { marked } from 'marked'
import { UploadFilled } from '@element-plus/icons-vue'
import { useUserStore } from '@/stores'
import { useRoute } from 'vue-router'
import { queryProblemDetailService, submitProblemService, queryResultService } from '@/api/problemController'
import { getRootCommentListService, addCommentService } from '@/api/commentController'
// Ace Editor
import { VAceEditor } from 'vue3-ace-editor'
import ace from 'ace-builds'

// 设置 Ace 资源路径
ace.config.set('basePath', '/ace')

// 基础扩展
import 'ace-builds/src-noconflict/ext-language_tools'

// 语言文件
import 'ace-builds/src-noconflict/mode-c_cpp'
import 'ace-builds/src-noconflict/mode-golang'
import 'ace-builds/src-noconflict/mode-python'
import 'ace-builds/src-noconflict/mode-java'

// 主题
import 'ace-builds/src-noconflict/theme-github'
import 'ace-builds/src-noconflict/theme-xcode'

// snippets
import 'ace-builds/src-noconflict/snippets/c_cpp'
import 'ace-builds/src-noconflict/snippets/golang'
import 'ace-builds/src-noconflict/snippets/java'
import 'ace-builds/src-noconflict/snippets/python'

const userStore = useUserStore()

// 题目详细信息
onMounted(() => {
    // 获取题目信息
    getProblemDetail()
})
const route = useRoute()
const problem = ref<API.Problem>({
    id: 0,
    title: '',
    description: '',
    level: 0,
    tags: [],
    status: 1
})
const getProblemDetail = async () => {
    const res = await queryProblemDetailService(Number(route.params.id))
    problem.value.id = Number(res.data.data.problem.id)
    problem.value.title = res.data.data.problem.title
    problem.value.description = res.data.data.problem.description
    problem.value.level = res.data.data.problem.level
    problem.value.tags = res.data.data.problem.tags
    problem.value.status = res.data.data.problem.status

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
const form = ref<API.SubmitForm>({
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
    form.value.lang = langMap[form.value.lang]
    const resp = await submitProblemService(form.value)
    loading.value = false

    // 轮训判题结果
    running.value = true
    timer = setInterval(() => {
        checkResult(resp.data.data.taskId)
    }, 1000)
}
const langMap: Record<string, string> = {
    c: 'c',
    c_cpp: 'cpp',
    java: 'java',
    python: 'python',
    golang: 'go',
}

// 轮训判题结果
let timer = 0
const MAX_RETRY = 20
let retryCount = 0
const running = ref(false)
const checkResult = async (task_id: string) => {
    retryCount++

    const resp = await queryResultService(task_id)
    if (resp.data.data !== null) {
        result_abstract.value = resp.data.data
        clearInterval(timer)
        running.value = false
        retryCount = 0
        return
    }

    // 超时判定
    if (retryCount >= MAX_RETRY) {
        clearInterval(timer)
        running.value = false
        retryCount = 0

        ElMessage({
            message: '请求超时，请稍后重试或查看记录',
            type: 'warning'
        })
    }
}
const result_abstract = ref<API.JudgeResultAbstract>()
const btype = computed(() => {
    if (result_abstract.value?.accepted) {
        return 'success'
    } else {
        return 'danger'
    }
})

// ACE主题和配置
ace.require('ace/ext/language_tools')

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

const theme = ref('xcode');

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

// 评论区
const place = computed(() => {
    if (userStore.userInfo.uid === 0) {
        return '请先登录...'
    } else {
        return '请输入评论...'
    }
})
const new_comment_text = ref('')
const onSubmitComment = async () => {
    const form = <API.AddCommentForm>{
        obj_id: problem.value.id,
        user_id: userStore.userInfo.uid,
        user_name: userStore.userInfo.username,
        user_avatar_url: userStore.userInfo.avatarUrl,
        content: new_comment_text.value
    }
    const res = await addCommentService(form)
    if (res.data.code === 0) {
        // 立刻渲染一条评论到第一条
        const new_comment = <API.Comment>({
            id: 0,
            objId: problem.value.id,
            userId: userStore.userInfo.uid,
            userName: userStore.userInfo.username,
            userAvatarUrl: userStore.userInfo.avatarUrl,
            content: new_comment_text.value,
            status: 1,
            replyCount: 0,
            likeCount: 0,
            childCount: 0,
            pubStamp: Date.now(),
            pubRegion: ' ',
            isRoot: true
        })
        root_comment_list.value.unshift(new_comment)
        new_comment_text.value = ''
    }
}

const root_comment_query_params = ref(<API.QueryRootCommentListParams>{
    page: 1,
    page_size: 5,
})
const root_comment_list = ref<API.Comment[]>([])
const root_comment_count = ref(0)
const getCommentList = async () => {
    root_comment_query_params.value.obj_id = problem.value.id

    const resp = await getRootCommentListService(root_comment_query_params.value)
    if (resp.data.data.data) {
        root_comment_list.value = resp.data.data.data
    }
    root_comment_count.value = Number(resp.data.data.total)
}
const handleCurrentChange = (page: number) => {
    root_comment_query_params.value.page = page
    getCommentList()
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
                        <el-tag v-if="problem.level === 1" type="success" round effect="plain">简单</el-tag>
                        <el-tag v-else-if="problem.level === 2" type="warning" round effect="plain">中等</el-tag>
                        <el-tag v-else type="danger" round effect="plain">困难</el-tag>
                    </div>
                    <div class="header-foot-tags">
                        <div v-for="item in problem.tags" :key="item">
                            <el-tag round effect="plain">{{ item }}</el-tag>
                        </div>
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
                <div class="result-area" v-if="result_abstract">
                    <el-tag :type="btype" size="large">
                        <template #default>
                            <strong>{{ result_abstract?.message }}</strong>
                        </template>
                    </el-tag>
                </div>

                <div class="submit-btn-area">
                    <el-button class="submit-btn" type="warning" :loading="loading" @click="onSubmit">
                        <el-icon>
                            <UploadFilled />
                        </el-icon>
                        &nbsp;Submit
                    </el-button>
                </div>
            </div>
        </div>

        <el-divider />

        <!-- 评论区 -->
        <div class="comment-area">
            <!-- 输入评论 -->
            <div class="comment-input">
                <el-input v-model="new_comment_text" type="textarea" show-word-limit :placeholder="place" resize="none"
                    :autosize="{ minRows: 5 }" maxlength="1000"
                    input-style="font-size: 18px; border: none; outline: none; overflow: hidden;" />
                <div style="display: flex; justify-content: flex-end;">
                    <el-button type="success" :disabled="!userStore.userInfo.uid || !new_comment_text" auto-insert-space
                        style="margin-top: 5px;" @click="onSubmitComment">评论</el-button>
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
                    <div class="comment-item" v-for="item in root_comment_list" :key="item.id">
                        <RootCommentItem :comment_data="item" :obj_id="problem.id"></RootCommentItem>
                    </div>
                    <!-- 分页 -->
                    <el-pagination v-model:current-page="root_comment_query_params.page"
                        v-model:page-size="root_comment_query_params.page_size" :total="root_comment_count"
                        :background="true" layout="prev, pager, next, jumper" @current-change="handleCurrentChange"
                        style="margin-top: 20px; justify-content: flex-end;" />
                </div>
            </el-card>
        </div>
    </el-card>
</template>

<style lang="less" scoped>
.container {
    width: 95%;
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
                display: flex;
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
            height: 400px;
            width: 100%;
            background-color: #ffffff;
            border: 1px solid #dcdfe6;
            padding: 8px;
            box-sizing: border-box;
        }

        /* 修改Ace编辑器行号区域背景颜色和宽度 */
        ::v-deep .ace_gutter {
            background-color: #f5f7fa;
        }
        
        /* 调整行号单元格的内边距，进一步减小行号区域的视觉宽度 */
        ::v-deep .ace_gutter-cell {
            padding-left: 4px;
            padding-right: 4px;
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
                flex: 1;
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

            .comment-item {
                margin-top: 15px;

                &:first-child {
                    margin-top: 0px;
                }

                margin-bottom: 15px;

                &:last-child {
                    margin-bottom: 0px;
                }
            }
        }
    }
}
</style>