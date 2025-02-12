<script lang="ts" setup>
import { ref, computed, defineProps } from 'vue'
// import type { Comment } from '@/api/interfaces/comment'

const props = defineProps<{
    comment_data: API.Comment
}>()

// 跳转链接
const user_href = computed(() => {
    return `/user/${comment_data.value.user_id}`
})

// 默认头像
const defaultAvatar = 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png'
</script>

<template>
    <div class="container">
        <!-- header区域 -->
        <div class="header-area">
            <!-- 头像 -->
            <a class="avatar-link" :href="user_href" target="_blank">
                <img :src="comment_data.user_avatar_url || circleUrl" class="avatar-image" />
            </a>
            <!-- 用户名及发布情况 -->
            <div class="user-pub">
                <div class="user-name">用户名</div>
                <div class="public-info">发布于 四川 2025.01.05</div>
            </div>
        </div>
        <!-- 内容区域 -->
        <div class="content-area">
            0123456789
        </div>
    </div>
</template>

<style lang="less" scoped>
.container {
    width: 100% !important;
    display: block;
    margin: 0px 0px !important;

    .header-area {
        width: 100%;
        display: flex;
        align-items: center;
        margin: auto;
        transition: background-color 0.3s ease;

        .avatar-link {
            display: flex;
            align-items: center;
            text-decoration: none;
            color: inherit;

            .avatar-image {
                width: 35px;
                height: 35px;
                border-radius: 50%;
                margin-right: 10px;
            }
        }

        .user-pub {
            width: 100%;
            display: flex;
            justify-content: space-between;
        }
    }

    .content-area {
        width: 100%;
        margin: auto;
        margin-top: 10px;
        border-style: solid;
        border-color: rgb(141, 0, 0);
    }
}
</style>