<script lang="ts" setup>
import { computed, ref } from 'vue'
import { formatTime } from '@/utils/format'
import { useUserStore } from '@/stores'
import { likeCommentService } from '@/api/commentController'
import { getChildCommentListService, addCommentService } from '@/api/commentController'

const userStore = useUserStore()

const props = defineProps<{
    comment_data: API.Comment
    obj_id: number,
}>()

// 跳转链接
const user_href = computed(() => {
    return `/user/${props.comment_data.user_id}`
})

// 默认头像
const defaultAvatar = 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png'

// root评论点赞
let is_already_like = ref(false)
const onClickLike = async () => {
    if (!userStore.userInfo.uid) {
        ElMessage.warning('请先登录')
        return
    }

    const form = <API.CommentLikeForm>{
        obj_id: props.obj_id,
        comment_id: props.comment_data.id,
    }
    const resp = await likeCommentService(form)
    if (resp.data.message === 'success') {
        props.comment_data.like_count++
        is_already_like.value = true
    }
}
// root评论回复
const show_replay_area = ref(false)
const onClickReply = () => {
    show_replay_area.value = !show_replay_area.value
    reply_comment_input.value = ''
}
const reply_comment_input = ref('')
const place = computed(() => {
    if (userStore.userInfo.uid === 0) {
        return '请先登录...'
    } else {
        return '请输入评论...'
    }
})
const onCancelReplay = () => {
    show_replay_area.value = false
    reply_comment_input.value = ''
}
const onSubmitComment = async () => {
    // 回复的是root评论，所以root_id和reply_id一样
    const form = <API.AddCommentForm>{
        obj_id: props.obj_id,
        user_id: userStore.userInfo.uid,
        user_name: userStore.userInfo.nickname,
        user_avatar_url: userStore.userInfo.avatar_url,
        content: reply_comment_input.value,
        is_root: false,
        root_id: props.comment_data.user_id,
        root_comment_id: props.comment_data.id,
        reply_id: props.comment_data.user_id,
        reply_comment_id: props.comment_data.id,
        reply_user_name: props.comment_data.user_name,
    }
    const res = await addCommentService(form)
    if (res.data.message === 'OK') {
        show_replay_area.value = false
        // 立刻渲染一条评论到第一条
        const new_comment = <API.Comment>{
            id: 0,
            obj_id: props.obj_id,
            user_id: userStore.userInfo.uid,
            user_name: userStore.userInfo.nickname,
            user_avatar_url: userStore.userInfo.avatar_url,
            content: reply_comment_input.value,
            status: 1,
            reply_count: 0,
            like_count: 0,
            child_count: 0,
            pub_stamp: Date.now(),
            pub_region: ' ',
            is_root: false,
            root_id: props.comment_data.user_id,
            root_comment_id: props.comment_data.id,
            reply_id: props.comment_data.user_id,
            reply_comment_id: props.comment_data.id,
        }
        child_list.value.unshift(new_comment)
        reply_comment_input.value = ''
    }
}

// 子评论回复
const child_comment_reply_action = async (form: API.AddCommentForm) => {
    const res = await addCommentService(form)
    if (res.data.message === 'success') {
        // 立刻渲染一条评论到第一条
        const new_comment = <API.Comment>{
            id: 0,
            obj_id: form.obj_id,
            user_id: form.user_id,
            user_name: form.user_name,
            user_avatar_url: form.user_avatar_url,
            content: form.content,
            status: 1,
            reply_count: 0,
            like_count: 0,
            child_count: 0,
            pub_stamp: Date.now(),
            pub_region: ' ',
            is_root: false,
            root_id: form.root_id,
            root_comment_id: form.root_comment_id,
            reply_id: form.reply_id,
            reply_comment_id: form.reply_comment_id,
            reply_user_name: form.reply_user_name
        }
        child_list.value.unshift(new_comment)
    }
}

// 展开回复列表
let is_expand = ref(false)
const expand_msg = computed(() => {
    return is_expand.value ? '收起' : '展开'
})
const onClickExpand = () => {
    is_expand.value = !is_expand.value
    reset_child()
    
    if (is_expand.value) {
        // 拉取子评论数据
        getChildCommentList()
    }
}
const child_list_cursor = ref(1)
const reset_child = () => {
    child_list_cursor.value = 1
    show_count.value = 0
    child_list.value = []
}
const child_list = ref(<API.Comment[]>[])
const child_total = ref(0)
const getChildCommentList = async () => {
    const params = <API.QueryChildCommentListParams>{
        obj_id: props.obj_id,
        root_id: props.comment_data.user_id,
        root_comment_id: props.comment_data.id,
        cursor: child_list_cursor.value
    }
    const resp = await getChildCommentListService(params)
    child_total.value = resp.data.data.total
    child_list_cursor.value = resp.data.data.cursor
    // 插入数组的条件：list不为nil，展示的数据量小于total
    if (resp.data.data.list && show_count.value < child_total.value) {
        resp.data.data.list.forEach((item: API.Comment) => {
            child_list.value.push(item)
            show_count.value++
        })
    }
}
const show_count = ref(0)
const need_show_more = computed(() => {
    return show_count.value < child_total.value
})
const onClickShowMore = () => {
    getChildCommentList()
}
const onClickHide = () => {
    is_expand.value = false
    reset_child()
}

</script>

<template>
    <div class="container">
        <!-- 顶层评论 -->
        <div class="root">
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
                        <svg v-if="!is_already_like" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"
                            width="1.2em" height="1.2em" fill="currentColor"
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
                <el-button plain @click="onClickExpand()" v-if="comment_data.reply_count > 0">
                    <template #default>
                        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="1.2em" height="1.2em"
                            fill="currentColor"
                            class="w-4.5 h-4.5 text-gray-6 dark:text-dark-gray-6 group-hover:text-gray-7 dark:group-hover:text-dark-gray-7">
                            <path fill-rule="evenodd"
                                d="M11.997 21.5a9.5 9.5 0 01-8.49-5.251A9.38 9.38 0 012.5 11.997V11.5c.267-4.88 4.12-8.733 8.945-8.999L12 2.5a9.378 9.378 0 014.25 1.007A9.498 9.498 0 0121.5 12a9.378 9.378 0 01-.856 3.937l.838 4.376a1 1 0 01-1.17 1.17l-4.376-.838a9.381 9.381 0 01-3.939.856zm3.99-2.882l3.254.623-.623-3.253a1 1 0 01.09-.64 7.381 7.381 0 00.792-3.346 7.5 7.5 0 00-4.147-6.708 7.385 7.385 0 00-3.35-.794H11.5c-3.752.208-6.792 3.248-7.002 7.055L4.5 12a7.387 7.387 0 00.794 3.353A7.5 7.5 0 0012 19.5a7.384 7.384 0 003.349-.793 1 1 0 01.639-.09z"
                                clip-rule="evenodd"></path>
                        </svg>
                        <span style="margin-left: 3px;">{{ expand_msg }} {{ comment_data.child_count }}条回复 </span>
                    </template>
                </el-button>
                <el-button plain @click="onClickReply()">
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
            <!-- 回复区域 -->
            <div class="reply-area" v-if="show_replay_area">
                <el-input v-model="reply_comment_input" type="textarea" show-word-limit :placeholder="place"
                    resize="none" :autosize="{ minRows: 5 }" maxlength="1000"
                    input-style="font-size: 18px; border: none; outline: none; overflow: hidden;" />
                <div style="display: flex; justify-content: flex-end;">
                    <el-button type="info" auto-insert-space style="margin-top: 5px;"
                        @click="onCancelReplay">取消</el-button>
                    <el-button type="success" :disabled="!userStore.userInfo.uid || !reply_comment_input"
                        auto-insert-space style="margin-top: 5px;" @click="onSubmitComment">评论</el-button>
                </div>
            </div>
        </div>
        <!-- 子评论区域 -->
        <div class="reply-list" v-if="is_expand">
            <SecondCommentItem v-for="item in child_list" :key="item.id" :comment_data="item" @child_comment_reply="child_comment_reply_action"></SecondCommentItem>
            <div class="replay-foot" v-if="need_show_more">
                <el-button plain @click="onClickShowMore">显示更多</el-button>
                <el-button plain @click="onClickHide">隐藏</el-button>
            </div>
        </div>
    </div>
</template>

<style lang="less" scoped>
.container {
    width: 100% !important;
    display: block;
    margin: 0px 0px !important;

    .root {
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

        .reply-area {}
    }

    .reply-list {
        margin-left: 15px;

        .SecondCommentItem {
            margin-top: 5px;
            margin-bottom: 5px;

            &:last-child {
                margin-bottom: 0px;
            }

            &:first-child {
                margin-top: 0px;
            }
        }

        .replay-foot {
            width: 100%;
            display: flex;
            justify-content: space-between;

            .el-button {
                border: none;
                padding: 0px 0px;
            }
        }
    }
}
</style>