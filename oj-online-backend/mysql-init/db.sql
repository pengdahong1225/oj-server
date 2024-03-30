use oj_online_server;

-- 用户信息表
create table if not exists user_info
(
    id BIGINT AUTO_INCREMENT,
    create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    delete_at TIMESTAMP,

    phone BIGINT not NULL comment '手机号',
    password VARCHAR(64) NOT null comment '密码',
    nickname VARCHAR(64) DEFAULT '新用户',
    email VARCHAR(64) DEFAULT '',
    gender tinyint DEFAULT 0 comment '0:woman 1:man',
    role tinyint DEFAULT 0 comment '0:user 1:admin',
    head_url VARCHAR(256) DEFAULT '' comment '头像url',

    pass_count BIGINT DEFAULT 0 comment '总题目AC数量',
    submit_count BIGINT DEFAULT 0 comment '总题目提交数量',

    PRIMARY KEY(id),
    UNIQUE INDEX idx_phone(phone)
)engine = InnoDB charset = utf8mb4;

-- 题目表
create table if not exists question
(
    id BIGINT AUTO_INCREMENT,
    create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    delete_at TIMESTAMP,

    title VARCHAR(64) NOT null comment '题目标题',
    level tinyint DEFAULT 0 comment '题目难度 0:简单 1:中等 2:困难',
    tags VARCHAR(64) DEFAULT '' comment '题目标签，#做前缀', -- '#数组 #双指针 #哈希表'

    description TEXT NOT null comment '题目描述', -- markdown
    test_case TEXT NOT null comment '测试用例', -- json{input,output}
    template TEXT NOT null comment '模板代码', -- json{lang,code}

    PRIMARY KEY(id)
)engine = InnoDB charset = utf8mb4;

-- 用户提交记录表，在线运行不用记录，提交代码需要记录
create table if not exists user_submit
(
    id BIGINT AUTO_INCREMENT,
    create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    delete_at TIMESTAMP,

    user_id BIGINT not null,
    question_id BIGINT not null,
    code TEXT NOT null comment '提交代码',
    result TEXT NOT null comment '运行结果',
    lang VARCHAR(64) DEFAULT '' comment '语言',

    PRIMARY KEY(id),
    FOREIGN KEY(user_id) REFERENCES user_info(id),
    FOREIGN KEY(question_id) REFERENCES question(id)
)engine = InnoDB charset = utf8mb4;
