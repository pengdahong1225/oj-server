<script lang="ts" setup>
import { ref, computed } from 'vue'
import { formatTime } from '@/utils/format'
import { useUserStore } from '@/stores'
import { likeCommentService } from '@/api/commentController'

const userStore = useUserStore()

const props = defineProps<{
    comment_data: API.Comment
}>()

// 跳转链接
const user_href = computed(() => {
    return `/user/${props.comment_data.user_id}`
})
// 默认头像
const defaultAvatar = 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png'

// 子评论点赞
let is_already_like = ref(false)
const onClickLike = async () => {
    if (!userStore.userInfo.uid) {
        ElMessage.warning('请先登录')
        return
    }

    const form = <API.CommentLikeForm>{
        obj_id: props.comment_data.obj_id,
        comment_id: props.comment_data.id,
    }
    const res = await likeCommentService(form)
    if (res.data.message === 'OK') {
        props.comment_data.like_count++
        is_already_like.value = true
    }
}
// 回复
const onClickReply = (id: number) => {

}
</script>

<template>
    <div class="container">
        <!-- header区域 -->
        <div class="header-area">
            <!-- 头像 -->
            <a class="avatar-link" :href="user_href" target="_blank">
                <img :src="comment_data.user_avatar_url || defaultAvatar" class="avatar-image" />
            </a>
            <!-- 用户名及发布情况 -->
            <div class="user-pub">
                <div class="user-name">{{ comment_data.user_name }}</div>
                <div class="public-info">发布于 {{ comment_data.pub_region || '未知地区' }} {{
                    formatTime(comment_data.pub_stamp) }}</div>
            </div>
        </div>
        <!-- 内容区域 -->
        <div class="content-area">
            {{ comment_data.content }}
        </div>
        <!-- 操作区域 -->
        <div class="operation-area">
            <el-button plain @click="onClickLike()">
                <template #default>
                    <svg v-if="!is_already_like" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="1.2em"
                        height="1.2em" fill="currentColor"
                        class="w-[18px] h-[18px] cursor-pointer text-gray-6 dark:text-dark-gray-6 group-hover:text-gray-7 dark:group-hover:text-dark-gray-7">
                        <path fill-rule="evenodd"
                            d="M7.04 9.11l3.297-7.419a1 1 0 01.914-.594 3.67 3.67 0 013.67 3.671V7.33h4.028a2.78 2.78 0 012.78 3.2l-1.228 8.01a2.778 2.778 0 01-2.769 2.363H5.019a2.78 2.78 0 01-2.78-2.78V11.89a2.78 2.78 0 012.78-2.78H7.04zm-2.02 2a.78.78 0 00-.781.78v6.232c0 .431.35.78.78.78H6.69V11.11H5.02zm12.723 7.793a.781.781 0 00.781-.666l1.228-8.01a.78.78 0 00-.791-.898h-5.04a1 1 0 01-1-1V4.77c0-.712-.444-1.32-1.07-1.56L8.69 10.322v8.58h9.053z"
                            clip-rule="evenodd"></path>
                    </svg>
                    <svg v-else xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="1.2em" height="1.2em"
                        fill="currentColor"
                        class="w-[18px] h-[18px] cursor-pointer text-green-s dark:text-dark-green-s hover:text-green-s dark:hover:text-dark-green-s">
                        <path fill-rule="evenodd"
                            d="M7.04 9.11l3.297-7.419a1 1 0 01.914-.594 3.67 3.67 0 013.67 3.671V7.33h4.028a2.78 2.78 0 012.78 3.2l-1.228 8.01a2.778 2.778 0 01-2.769 2.363H5.019a2.78 2.78 0 01-2.78-2.78V11.89a2.78 2.78 0 012.78-2.78H7.04zm-2.02 2a.78.78 0 00-.781.78v6.232c0 .431.35.78.78.78H6.69V11.11H5.02z"
                            clip-rule="evenodd"></path>
                    </svg>
                    <span style="margin-left: 3px;"> {{ comment_data.like_count }} </span>
                </template>
            </el-button>
            <el-button plain @click="onClickReply(comment_data.id)">
                <template #default>
                    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" width="1.2em" height="1.2em"
                        fill="currentColor"
                        class="w-4.5 h-4.5 text-gray-6 dark:text-dark-gray-6 group-hover:text-gray-7 dark:group-hover:text-dark-gray-7">
                        <path fill-rule="evenodd"
                            d="M5.83 2.106c.628-.634 1.71-.189 1.71.704v2.065c4.821.94 6.97 4.547 7.73 8.085l-.651.14.652-.134c.157.757-.83 1.192-1.284.565l-.007-.009c-1.528-2.055-3.576-3.332-6.44-3.502v2.352c0 .893-1.082 1.338-1.71.704L1.091 8.295a1 1 0 010-1.408l4.737-4.78zm7.303 8.617C12.08 8.495 10.204 6.68 7.046 6.14c-.47-.08-.84-.486-.84-.99V3.62L2.271 7.591l3.934 3.971V9.667a.993.993 0 011.018-.995c2.397.065 4.339.803 5.909 2.051z"
                            clip-rule="evenodd"></path>
                    </svg>
                    <span style="margin-left: 3px;"> 回复 </span>
                </template>
            </el-button>
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
    }

    .operation-area {
        .el-button {
            border: none;
            padding: 0px 0px;
            margin-left: 5px;
            margin-right: 5px;
        }
    }
}
</style>