use oj_online_server;
-- 题目表
create table if not exists problem
(
    id BIGINT AUTO_INCREMENT,
    create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    delete_at TIMESTAMP,

    title VARCHAR(64) NOT null comment '题目标题',
    level tinyint DEFAULT 0 comment '题目难度 1:简单 2:中等 3:困难',
    tags JSON comment '题目标签',
    description TEXT NOT null comment '题目描述',
    create_by BIGINT DEFAULT 0 comment '题目创建者',
    comment_count BIGINT DEFAULT 0 comment '评论总数量',

    state tinyint DEFAULT 1 comment '状态 1：发布 0：隐藏',

    PRIMARY KEY(id),
    UNIQUE INDEX idx_title(title)
)engine = InnoDB charset = utf8mb4;
