<script lang="ts" setup>
import { ref } from 'vue'
import type { QueryRecordListParams } from '@/types/record'

const loading = ref(false)
const params = ref<QueryRecordListParams>({
    page: 1,
    page_size: 10
})
const 

</script>

<template>
    <el-card class="status-list" shadow="hover">
        <!-- 表单区域 -->
        <!-- <el-form inline class="form">
            <el-form-item style="margin-right: 15px;">
                <el-input v-model="params.keyword" style="width: 240px" placeholder="problem name"
                    :suffix-icon="Search" />
            </el-form-item>
            <el-form-item>
                <el-button @click="onSearch" type="primary">搜索</el-button>
                <el-button @click="onReset">重置</el-button>
            </el-form-item>
        </el-form> -->

        <!-- 表格区域 -->
        <el-table v-loading="loading" :data="problemList">
            <el-table-column label="Status" prop="status" width="80" align="center">
                <template #default="{ row }">
                    <el-icon v-if="row.status === 1" color="green" size="18"><Select /></el-icon>
                </template>
            </el-table-column>

            <el-table-column label="#" prop="id" width="80">
                <template #default="{ row }">
                    <el-link type="primary" :underline="false" @click="
                        $router.push({
                            path: `/problem/${row.id}`
                        })
                        ">{{ row.id }}</el-link>
                </template>
            </el-table-column>

            <el-table-column label="Title" prop="title">
                <template #default="{ row }">
                    <el-link type="primary" :underline="false" @click="
                        $router.push({
                            path: `/problem/${row.id}`
                        })
                        ">{{ row.title }}</el-link>
                </template>
            </el-table-column>
            <el-table-column label="Level" prop="level">
                <template #default="{ row }">
                    <el-tag v-if="row.level === 1" type="primary">简单</el-tag>
                    <el-tag v-else-if="row.level === 2" type="warning">中等</el-tag>
                    <el-tag v-else type="danger">困难</el-tag>
                </template>
            </el-table-column>
            <el-table-column label="Tags">
                <template #default="{ row }">
                    <el-tag v-for="tag in row.tags" :key="tag" style="margin-left: 3px;margin-right: 3px;">{{ tag
                        }}</el-tag>
                </template>
            </el-table-column>

            <template #empty>
                <el-empty description="没有数据"></el-empty>
            </template>
        </el-table>

        <!-- 分页 -->
        <el-pagination v-model:current-page="params.page" v-model:page-size="params.page_size" :total="total"
            :background="true" layout="prev, pager, next, jumper" @current-change="handleCurrentChange"
            style="margin-top: 20px; justify-content: flex-end;" />
    </el-card>
</template>

<style lang="scss" scoped></style>