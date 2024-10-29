<script lang="ts" setup>
import { ref } from 'vue'
// import default_banner from '@/assets/default_banner.jpg'
import { formatTime } from '@/utils/format'

const loading = ref(false)
const announcement_list = ref([
    {
        id: 1,
        title: 'frontend',
        stamp: '1730172635',
    },
    {
        id: 2,
        title: 'backend',
        stamp: '1730172198',
    },
])
const total = ref(2)
const params = ref({
    page: 1,
    page_size: 10, // page_size默认为10
    keyword: '',
})
const handleCurrentChange = (page: number) => {
    params.value.page = page
    console.log(params.value)
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

                <el-table v-loading="loading" :data="announcement_list">
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
                    <el-table-column label="Time" prop="stamp">
                        <template #default="{ row }">
                            {{ formatTime(row.stamp) }}
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

            </el-card>
        </div>

    </div>
</template>

<style lang="scss" scoped>
.container {
    display: block;
    width: 60%;
    margin: auto;

    .banner {
        border: 1px solid #00c853;
        margin-bottom: 20px;
        border-radius: 8px;
    }

    .notice {
        display: flex;
        border-radius: 8px;
        align-items: flex-start; /* 确保卡片顶部对齐，而不会拉伸高度 */

        .announcements {
            width: 70%;
            margin-right: 6px;
        }

        .search {
            width: 30%;
            margin-left: 6px;
        }
    }
}
</style>