<script lang="ts" setup>
import { useRoute } from 'vue-router'
import { queryRecordService } from '@/api/userController'
import { onMounted, ref } from 'vue'

onMounted(() => {
    getRecord()
})

const route = useRoute();

const id = route.params.id
const record = ref(<API.Record>({}))
const getRecord = async () => {
    const res = await queryRecordService(Number(id))
    console.log(res)
    record.value = res.data.data
    console.log(record.value)
}

</script>

<template>
    <div class="container">
        <div class="status-header">
            <div class="state">
                <span>通过</span>
                <span>63 / 63 个通过的测试用例</span>
            </div>
            <div class="submit-info">黑手双城 提交于 2024.02.28 22:42</div>
        </div>
        <div class="status-detail">
            详细信息
        </div>
        <div class="status-code">
            <code>
                <pre>{{ record.code }}</pre>
            </code>
        </div>
    </div>
</template>

<style lang="less" scoped>
.container {
    display: block;
    margin: auto;
    width: 70%;

    .status-header {
        display: block;
    }
}
</style>