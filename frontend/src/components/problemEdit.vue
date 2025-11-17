<script setup lang="ts">
import { ref } from 'vue'
import {
    addProblemService,
    updateProblemService,
    queryProblemDetailService
} from '@/api/problemController'

// 父子组件通信数据
const open = (mode: string, problem_id: number) => {
    drawerVisible.value = true
    edit_mode.value = mode
    if (mode === 'update') {
        getProblemDetail(problem_id)
    }
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
    if (resp.data.code === 0) {
        ElMessage.success('提交成功')
    } else {
        ElMessage.error(resp.data.message || '提交失败')
    }

    drawerVisible.value = false
    emit('success', edit_mode.value)
}

const level_list = [
    { label: '简单', value: 1 },
    { label: '中等', value: 2 },
    { label: '困难', value: 3 }
]
const tag_list = [
  "数组",
  "字符串",
  "哈希表",
  "动态规划",
  "数学",
  "排序",
  "贪心",
  "深度优先搜索",
  "二分查找",
  "数据库",
  "树",
  "矩阵",
  "广度优先搜索",
  "位运算",
  "双指针",
  "前缀和",
  "二叉树",
  "堆（优先队列）",
  "模拟",
  "栈",
  "图",
  "计数",
  "滑动窗口",
  "设计",
  "回溯",
  "枚举",
  "链表",
  "并查集",
  "数论",
  "有序集合",
  "单调栈",
  "线段树",
  "分治",
  "字典树",
  "递归",
  "组合数学",
  "状态压缩",
  "队列",
  "二叉搜索树",
  "几何",
  "记忆化搜索",
  "树状数组",
  "哈希函数",
  "拓扑排序",
  "最短路",
  "字符串匹配",
  "滚动哈希",
  "博弈",
  "数据流",
  "交互",
  "单调队列",
  "脑筋急转弯",
  "双向链表",
  "归并排序",
  "随机化",
  "计数排序",
  "快速选择",
  "迭代器",
  "概率与统计",
  "多线程",
  "扫描线",
  "后缀数组",
  "桶排序",
  "最小生成树",
  "Shell",
  "水塘抽样",
  "强连通分量",
  "欧拉回路",
  "基数排序",
  "双连通分量",
  "拒绝采样"
]

const drawerVisible = ref(false)
const formRef = ref()

// 默认创建
let edit_mode = ref('create')
const create_form = ref<API.CreateProblemForm>({
    problem_title: '',
    level: level_list[0]?.value ?? null,
    tags: [],
    description: '',
})
const update_form = ref<API.UpdateProblemForm>({
    problem_id:  0,
    problem_title: '',
    level: level_list[0]?.value ?? null,
    tags: [],
    description: '',
})

const getProblemDetail = async (id: number) => {
    const resp = await queryProblemDetailService(id)
    console.log(resp)
    update_form.value.problem_id = resp.data.data.problem_id
    update_form.value.problem_title = resp.data.data.problem_title
    update_form.value.description = resp.data.data.description
    update_form.value.level = resp.data.data.level
    update_form.value.tags = resp.data.data.tags
}

</script>

<template>
    <el-drawer v-model="drawerVisible" size="45%">
        <template #header>
            <strong v-if="edit_mode == 'create'"> 新增题目 </strong>
            <strong v-else> 编辑题目 </strong>
        </template>

        <template #default>
            <el-form v-if="edit_mode == 'create'" class="form" ref="formRef" :model="create_form" label-width="150px" size="large">
                <el-form-item label="标题" prop="title">
                    <el-input v-model="create_form.problem_title" placeholder="请输入题目标题" style="width: 50%"></el-input>
                </el-form-item>
                <el-form-item label="Level" prop="level">
                    <el-select size="large" v-model="create_form.level" placeholder="Select" style="width: 240px">
                        <el-option v-for="item in level_list" :key="item.value" :label="item.label"
                            :value="item.value" />
                    </el-select>
                </el-form-item>
                <el-form-item label="标签" prop="tags">
                    <el-select v-model="create_form.tags" multiple filterable allow-create default-first-option
                        placeholder="Select" style="width: 240px">
                        <el-option v-for="item in tag_list" :key="item" :label="item" :value="item" />
                    </el-select>
                </el-form-item>
                <el-form-item label="题目详情" prop="description">
                    <el-input type="textarea" :rows="10" resize="none" v-model="create_form.description"
                        placeholder="请描述题目详情" style="width: 100%"></el-input>
                </el-form-item>
                <el-form-item style="">
                    <el-button @click="drawerVisible = false">取消</el-button>
                    <el-button type="primary" @click="onSubmit"> 确认 </el-button>
                </el-form-item>
            </el-form>
            <el-form v-else class="form" ref="formRef" :model="update_form" label-width="150px" size="large">
                <el-form-item label="标题" prop="title">
                    <el-input v-model="update_form.problem_title" placeholder="请输入题目标题" style="width: 50%"></el-input>
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