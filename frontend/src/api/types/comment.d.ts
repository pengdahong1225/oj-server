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

    // 评论列表查询参数
    type QueryCommentListParams = {
        obj_id: number;
        root_id: number;
        root_comment_id: number;
        reply_id: number;
        reply_comment_id: number;

        cursor: number; // 分页参数
    }
}