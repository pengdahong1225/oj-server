namespace API {
    type Comment = {
        id: number;

        obj_id: number;
        user_id: number;
        user_name: string;
        user_avatar_url: string;
        content: string;
        status: number;
        reply_count: number;
        like_count: number;
        child_count: number;
        pub_stamp: number;
        pub_region: string; 

        is_root: boolean;
        root_id: number;
        root_comment_id: number;
        reply_id: number;
        reply_comment_id: number;
    }

    // 顶层评论列表查询参数
    type QueryRootCommentListParams = {
        obj_id: number;

        page: number;
        page_size: number;
    }
    // 二层评论列表查询参数
    type QueryChildCommentListParams = {
        obj_id: number;
        root_id: number;
        root_comment_id: number;
        reply_id: number;
        reply_comment_id: number;

        cursor: number;
    }

    // 提交评论接口参数
    type AddCommentForm = {
        obj_id: number;
        user_id: number;
        user_name: string;
        user_avatar_url: string;
        content: string;

        // 非顶层评论需要以下字段
        root_id: number;
        root_comment_id: number;
        reply_id: number;
        reply_comment_id: number;
    }
    // 评论点赞接口表单
    type CommentLikeForm = {
        obj_id: number;
        comment_id: number;
    }
}