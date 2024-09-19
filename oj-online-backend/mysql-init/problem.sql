-- 题目表
create table if not exists problem
(
    id BIGINT AUTO_INCREMENT,
    create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    delete_at TIMESTAMP,

    title VARCHAR(64) NOT null comment '题目标题',
    level tinyint DEFAULT 0 comment '题目难度 1:简单 2:中等 3:困难',
    tags VARCHAR(64) DEFAULT '' comment '题目标签，#做前缀', -- '#数组 #双指针 #哈希表'
    description TEXT NOT null comment '题目描述',
    create_by BIGINT DEFAULT 0 comment '题目创建者',
    comment_count BIGINT DEFAULT 0 comment '评论总数量',

    -- 三项配置文本格式：json
    test_case TEXT NOT null comment '测试用例',
    compile_config TEXT NOT null comment '编译配置',
    run_config TEXT NOT null comment '运行配置',

    PRIMARY KEY(id)
)engine = InnoDB charset = utf8mb4;
