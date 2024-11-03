<script lang="ts" setup>
import { onMounted, ref } from 'vue'
import default_banner from '@/assets/default_banner.jpg'
import { Search } from '@element-plus/icons-vue'
import { formatTime } from '@/utils/format'
import { queryNoticeListService } from '@/api/notice'
import type { QueryNoticeListParams, Notice } from '@/types/notice'

onMounted(() => {
    queryNoticeList()
})

const loading = ref(false)
const notice_list = ref(<Notice[]>[])
const total = ref(0)
const params = ref(<QueryNoticeListParams>{
    page: 1,
    page_size: 10, // page_size默认为10
    keyword: '',
})
const handleCurrentChange = (page: number) => {
    params.value.page = page
    console.log(params.value)
}
const handleSearch = () => {
    params.value.page = 1
    queryNoticeList()
}
const handleReset = () => {
    params.value.keyword = ''
    params.value.page = 1
    queryNoticeList()
}

const queryNoticeList = async () => {
    loading.value = true
    const res = await queryNoticeListService(params.value)
    console.log(res)
    notice_list.value = res.data.data.noticeList
    total.value = res.data.data.total
    loading.value = false
}

</script>

<template>
    <div class="container">
        <!-- 轮播图区域 -->
        <div class="banner">
            <el-carousel height="400px">
                <el-carousel-item v-for="item in 1" :key="item">
                    <img :src="default_banner" alt="" height="400px" width="100%" />
                </el-carousel-item>
            </el-carousel>
        </div>

        <!-- 公告区域 -->
        <div class="notice">
            <el-card class="announcements" shadow="hover">
                <template #header>
                    <strong>Announcements</strong>
                </template>

                <el-table v-loading="loading" :data="notice_list">
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
                    <el-table-column label="Time" prop="create_at">
                        <template #default="{ row }">
                            {{ formatTime(row.create_at) }}
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

            <el-card class="search" shadow="hover">
                <template #header>
                    <strong>Search</strong>
                </template>
                <div class="search-input">
                    <el-input v-model="params.keyword" placeholder="key word"
                        :prefix-icon="Search" />
                    <el-button type="primary" style="margin-left: 10px;" @click="handleSearch">搜索</el-button>
                    <el-button style="margin-left: 5px;" @click="handleReset">重置</el-button>
                </div>
            </el-card>
        </div>

    </div>
</template>

<style lang="scss" scoped>
.container {
    display: block;
    width: 70%;
    margin: auto;

    .banner {
        border: 1px solid #00c853;
        margin-bottom: 20px;
        border-radius: 8px;
    }

    .notice {
        display: flex;
        border-radius: 8px;
        align-items: flex-start;
        /* 确保卡片顶部对齐，而不会拉伸高度 */

        .announcements {
            width: 70%;
            margin-right: 6px;
        }

        .search {
            width: 30%;
            margin-left: 6px;
            .search-input{
                width: 100%;
                display: flex;
            }
        }
    }
}
</style>