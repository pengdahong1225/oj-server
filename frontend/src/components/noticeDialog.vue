<script lang="ts" setup>
import { ref } from 'vue'
import type { Notice } from '@/types/notice'

const dialogVisible = ref(false)
const notice = ref<Notice>({
    content: '',
    create_at: '',
    create_by: 0,
    id: 0,
    status: 0,
    title: ''
})
const open = (data: Notice) => {
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