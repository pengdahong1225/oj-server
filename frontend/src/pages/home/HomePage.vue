<script lang="ts" setup>
import { onMounted, ref } from 'vue'
import { DataBoard, Refresh, DataAnalysis } from '@element-plus/icons-vue'
import { formatTime } from '@/utils/format'
import { queryNoticeListService } from '@/api/noticeController'

onMounted(() => {
    queryNoticeList()
})

const loading = ref(false)
const notice_list = ref(<API.Notice[]>[])
const total = ref(0)
const params = ref(<API.QueryNoticeListParams>{
    page: 1,
    page_size: 10, // page_size默认为10
    keyword: '',
})
const handleCurrentChange = (page: number) => {
    params.value.page = page
    queryNoticeList() // 分页切换时重新加载数据
}

// 最近一周过题榜单（后端 leaderboard 尚未实现，暂空）
const rank_list = ref<Array<{ user_name: string; ac_total: number }>>([])

const queryNoticeList = async () => {
    loading.value = true
    try {
        const res = await queryNoticeListService(params.value)
        console.log(res)

        if (res.data && res.data.data) {
            notice_list.value = res.data.data.noticeList || []
            total.value = res.data.data.total || 0
        }
    } catch (error) {
        console.error('Failed to query notice list:', error)
        notice_list.value = []
        total.value = 0
    } finally {
        loading.value = false
    }
}

const dialogRef = ref()
const onClick = (index: number) => {
    console.log(index)
    dialogRef.value.open(notice_list.value[index])
}

// banner
const bannerList = [
    '',
]

</script>

<template>
    <div class="container">
        <!-- 轮播图区域 -->
        <el-card class="banner" shadow="hover">
            <template #header>
                <div class="banner-title">欢迎使用 OpenOJ</div>
            </template>
            <el-carousel height="600px">
                <el-carousel-item v-for="item in bannerList" :key="item">
                    <img :src="item" alt="" height="600px" width="100%" />
                </el-carousel-item>
            </el-carousel>
        </el-card>

        <!-- 公告区域 -->
        <div class="notice">
            <el-card class="announcements" shadow="hover">
                <template #header>
                    <div class="card-header">
                        <div class="card-header-left">
                            <el-icon><DataBoard /></el-icon>
                            <span>公告</span>
                        </div>
                        <el-button type="primary" :icon="Refresh" @click="queryNoticeList">刷新</el-button>
                    </div>
                </template>

                <el-table v-loading="loading" :data="notice_list">
                    <el-table-column label="#" prop="id" width="80">
                        <template #default="{ $index, row }">
                            <el-link type="primary" :underline="false" @click="onClick($index)">{{ row.id }}</el-link>
                        </template>
                    </el-table-column>
                    <el-table-column label="Title" prop="title">
                        <template #default="{ $index, row }">
                            <el-link type="primary" :underline="false" @click="onClick($index)">{{ row.title }}</el-link>
                        </template>
                    </el-table-column>
                    <el-table-column label="Time" prop="create_at">
                        <template #default="{ row }">
                            {{ formatTime(row.createAt) }}
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

            <el-card class="rank" shadow="hover">
                <template #header>
                    <div class="card-header">
                        <div class="card-header-left">
                            <el-icon><DataAnalysis /></el-icon>
                            <span>最近一周过题榜单</span>
                        </div>
                    </div>
                </template>
                <el-table :data="rank_list" stripe>
                    <el-table-column label="#" type="index" width="60" align="center" />
                    <el-table-column label="用户名" prop="user_name" align="center" />
                    <el-table-column label="通过" prop="ac_total" align="center" />
                    <template #empty>
                        <el-empty description="暂无数据"></el-empty>
                    </template>
                </el-table>
            </el-card>
        </div>

        <notice-dialog ref="dialogRef"></notice-dialog>
    </div>
</template>

<style lang="less" scoped>
.container {
    display: block;
    width: 80%;
    margin: auto;

    .banner {
        margin-bottom: 20px;

        .banner-title {
            text-align: center;
            font-weight: bold;
            font-size: 24px;
        }
    }

    .notice {
        display: flex;
        border-radius: 8px;
        align-items: flex-start;
        /* 确保卡片顶部对齐，而不会拉伸高度 */

        .card-header {
            display: flex;
            align-items: center;
            justify-content: space-between;

            .card-header-left {
                display: flex;
                align-items: center;
                gap: 6px;
                font-size: 18px;
                line-height: 1.5;
            }
        }

        .announcements {
            width: 65%;
            margin-right: 6px;
        }

        .rank {
            width: 35%;
            margin-left: 6px;
        }
    }
}

// 深度选择器，将链接颜色改为黑色
:deep(.el-link--primary) {
    color: #000 !important;
    
    &:hover {
        color: #333 !important;
    }
}
</style>