use oj_online_server;

-- 用户信息表
create table if not exists user_info
(
    id BIGINT AUTO_INCREMENT,
    create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    delete_at TIMESTAMP,

    mobile BIGINT not NULL comment '手机号',
    nickname VARCHAR(64) DEFAULT '新用户',
    email VARCHAR(64) DEFAULT '',
    gender tinyint DEFAULT 0 comment '0:woman 1:man',
    role tinyint DEFAULT 0 comment '0:user 1:admin',
    avatar_url VARCHAR(256) DEFAULT '' comment '头像url',

    PRIMARY KEY(id),
    UNIQUE INDEX idx_mobile(mobile)
)engine = InnoDB charset = utf8mb4;

-- 做题信息表
create table if not exists user_problem_statistics(
    id BIGINT AUTO_INCREMENT,
    create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    delete_at TIMESTAMP,
    
    uid BIGINT not null,
    submit_count BIGINT DEFAULT 0 comment '题目提交数量',
    accomplish_count BIGINT DEFAULT 0 comment '题目通过数量',
    easy_problem_count BIGINT DEFAULT 0 comment '通过的简单题目数量',
    medium_problem_count BIGINT DEFAULT 0 comment '通过的中等题目数量',
    hard_problem_count BIGINT DEFAULT 0 comment '通过的困难题目数量',

    PRIMARY KEY(id),
    FOREIGN KEY(uid) REFERENCES user_info(id),
    UNIQUE INDEX idx_uid(uid)
)engine = InnoDB charset = utf8mb4;

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

    -- 三项配置文本格式：json
    test_case TEXT NOT null comment '测试用例',
    compile_config TEXT NOT null comment '编译配置',
    run_config TEXT NOT null comment '运行配置',

    PRIMARY KEY(id)
)engine = InnoDB charset = utf8mb4;

-- 用户提交记录表
create table if not exists user_submit_record
(
    id BIGINT AUTO_INCREMENT,
    create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    delete_at TIMESTAMP,

    uid BIGINT not null,
    problem_id BIGINT not null,
    code TEXT NOT null comment '提交的代码',
    result TEXT NOT null comment '运行结果集', -- json
    lang VARCHAR(64) DEFAULT '' comment '语言',

    PRIMARY KEY(id),
    FOREIGN KEY(uid) REFERENCES user_info(id),
    FOREIGN KEY(problem_id) REFERENCES problem(id)
)engine = InnoDB charset = utf8mb4;

-- 用户解题表
create table if not exists user_solution
(
    id BIGINT AUTO_INCREMENT,
    create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    delete_at TIMESTAMP,

    uid BIGINT NOT NULL,
    problem_id BIGINT NOT NULL,

    PRIMARY KEY(id),
    FOREIGN KEY(uid) REFERENCES user_info(id),
    FOREIGN KEY(problem_id) REFERENCES problem(id),
    INDEX idx_uid(uid)
)engine = InnoDB charset = utf8mb4;