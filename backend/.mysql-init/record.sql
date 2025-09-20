use oj_online_server;

-- 用户提交记录表
create table if not exists user_submit_record
(
    id BIGINT AUTO_INCREMENT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    uid BIGINT not null comment '用户id',
    user_name VARCHAR(64) DEFAULT '' comment '用户名',
    problem_id BIGINT not null comment '题目id',
    problem_name VARCHAR(64) not null comment '题目名称',
    status VARCHAR(64) not null comment '状态',
    code TEXT NOT null comment '提交的代码',
    result blob NOT null comment '运行结果集',
    lang VARCHAR(64) DEFAULT '' comment '语言',

    PRIMARY KEY(id),
    INDEX idx_uid(uid)
)engine = InnoDB charset = utf8mb4;

-- 用户解题表
create table if not exists user_solution
(
    id BIGINT AUTO_INCREMENT,
    uid BIGINT NOT NULL COMMENT '用户id',
    problem_id BIGINT NOT NULL comment '题目id',
    level tinyint DEFAULT 0 comment '题目难度 1:简单 2:中等 3:困难',

    create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    delete_at TIMESTAMP,

    PRIMARY KEY(id),
    INDEX idx_uid(uid, problem_id)
)engine = InnoDB charset = utf8mb4;

-- 用户解题统计表 分区
create table if not exists user_problem_statistics(
    uid BIGINT not null comment '用户id',
    create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    delete_at TIMESTAMP,

    period CHAR(7) NOT NULL COMMENT 'YYYY-MM',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    submit_count INT DEFAULT 0 comment '题目提交数量',
    accomplish_count INT DEFAULT 0 comment '题目通过数量',
    easy_problem_count INT DEFAULT 0 comment '通过的简单题目数量',
    medium_problem_count INT DEFAULT 0 comment '通过的中等题目数量',
    hard_problem_count INT DEFAULT 0 comment '通过的困难题目数量',

    PRIMARY KEY(uid, period) -- 复合主键
)engine = InnoDB charset = utf8mb4;