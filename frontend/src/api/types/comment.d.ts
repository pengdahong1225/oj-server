namespace API {
    // 评论（响应 Comment，camelCase）
    type Comment = {
        id: number;
        objId: number;
        userId: number;
        userName: string;
        userAvatarUrl: string;
        content: string;
        status: number;
        replyCount: number;
        likeCount: number;
        childCount: number;
        pubStamp: number;
        pubRegion: string;

        isRoot?: boolean;
        rootId: number;
        rootCommentId: number;
        replyId: number;
        replyCommentId: number;
        replyUserName: string;
    }

    // 顶层评论列表查询参数（请求 query，snake_case 不变）
    type QueryRootCommentListParams = {
        obj_id: number;
        page: number;
        page_size: number;
    }
    // 二层评论列表查询参数（请求 query）
    type QueryChildCommentListParams = {
        obj_id: number;
        root_id: number;
        root_comment_id: number;
        reply_id: number;
        reply_comment_id: number;
        cursor: number;
    }

    // 提交评论表单（请求体，snake_case 不变）
    type AddCommentForm = {
        obj_id: number;
        user_id: number;
        user_name: string;
        user_avatar_url: string;
        content: string;

        // 非顶层评论需要以下字段
        root_id?: number;
        root_comment_id?: number;
        reply_id?: number;
        reply_comment_id?: number;
        reply_user_name?: string;
        is_root?: boolean;
    }
    // 评论点赞表单（请求体）
    type CommentLikeForm = {
        obj_id: number;
        comment_id: number;
    }
}
