<script lang="ts" setup>
import { onMounted, ref } from 'vue'
import { DataBoard, Refresh, DataAnalysis, Edit, User, Back } from '@element-plus/icons-vue'
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
        if (res.data && res.data.data) {
            total.value = res.data.data.total || 0
            notice_list.value = res.data.data.list || []
        }
    } catch (error) {
        console.error('failed to query notice list:', error)
        notice_list.value = []
        total.value = 0
    } finally {
        loading.value = false
    }
}

// 当前查看的公告（null 表示列表态）
const current_notice = ref<API.Notice | null>(null)
const openNotice = (item: API.Notice) => {
    current_notice.value = item
}
// 头部按钮：详情态返回列表 / 列表态刷新
const onHeaderButtonClick = () => {
    if (current_notice.value) {
        current_notice.value = null
    } else {
        queryNoticeList()
    }
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
                            <span>{{ current_notice ? current_notice.title : '公告' }}</span>
                        </div>
                        <el-button type="primary" :icon="current_notice ? Back : Refresh" @click="onHeaderButtonClick">
                            {{ current_notice ? '返回' : '刷新' }}
                        </el-button>
                    </div>
                </template>

                <!-- 公告详情 -->
                <div v-if="current_notice" class="notice-detail">
                    {{ current_notice.content }}
                </div>

                <!-- 公告列表 -->
                <div v-else v-loading="loading">
                    <div class="notice-list">
                        <div v-for="item in notice_list" :key="item.id" class="notice-item" @click="openNotice(item)">
                            <div class="notice-item-title">{{ item.title }}</div>
                            <div class="notice-item-meta">
                                <span class="meta-item">
                                    <el-icon><Edit /></el-icon>
                                    {{ formatTime(item.createAt) }}
                                </span>
                                <span class="meta-item">
                                    <el-icon><User /></el-icon>
                                    {{ item.createByName || item.createBy }}
                                </span>
                            </div>
                        </div>
                        <el-empty v-if="notice_list.length === 0" description="没有数据"></el-empty>
                    </div>
                    <!-- 分页 -->
                    <el-pagination v-model:current-page="params.page" v-model:page-size="params.page_size" :total="total"
                        :background="true" layout="prev, pager, next, jumper" @current-change="handleCurrentChange"
                        style="margin-top: 20px; justify-content: flex-end;" />
                </div>
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

            .notice-detail {
                min-height: 200px;
                padding: 12px 4px;
                white-space: pre-wrap;
                word-break: break-word;
                color: #303133;
                line-height: 1.7;
            }

            .notice-list {
                .notice-item {
                    border-left: 3px solid var(--el-color-primary);
                    padding: 8px 16px;
                    margin-bottom: 12px;
                    cursor: pointer;
                    background-color: #f8f9fb;
                    border-radius: 0 4px 4px 0;
                    transition: background-color 0.2s;

                    &:hover {
                        background-color: #eef1f6;
                    }

                    .notice-item-title {
                        font-size: 15px;
                        color: var(--el-color-primary);
                        margin-bottom: 6px;

                        &:hover {
                            text-decoration: underline;
                        }
                    }

                    .notice-item-meta {
                        display: flex;
                        align-items: center;
                        gap: 16px;
                        font-size: 13px;
                        color: #909399;

                        .meta-item {
                            display: inline-flex;
                            align-items: center;
                            gap: 4px;
                        }
                    }
                }
            }
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