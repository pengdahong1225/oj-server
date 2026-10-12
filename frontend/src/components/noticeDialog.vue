<script lang="ts" setup>
import { ref } from 'vue'

const dialogVisible = ref(false)
const notice = ref<API.Notice>({
    content: '',
    createAt: '',
    createBy: 0,
    id: 0,
    status: 0,
    title: ''
})
const open = (data: API.Notice) => {
    dialogVisible.value = true
    notice.value = data
}
defineExpose({
    open
})

</script>

<template>
    <el-dialog v-model="dialogVisible" :title="notice.title" width="70%">
        {{ notice.content }}
    </el-dialog>
</template>