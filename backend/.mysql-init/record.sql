use oj_online_server;

-- 用户提交记录表
create table if not exists user_submit_record
(
    id BIGINT AUTO_INCREMENT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    uid BIGINT not null comment '用户id',
    problem_id BIGINT not null comment '题目id',
    accepted boolean DEFAULT false comment '是否通过',
    message  VARCHAR(256) DEFAULT '' comment '测评结果描述',
    code TEXT NOT null comment '提交的代码',
    result blob NOT null comment '运行结果集',
    lang VARCHAR(64) DEFAULT '' comment '语言',

    PRIMARY KEY(id),
    UNIQUE INDEX uk_user_problem(uid, problem_id)
)engine = InnoDB charset = utf8mb4;

-- 用户解题表
create table if not exists user_solution
(
    id BIGINT AUTO_INCREMENT,
    uid BIGINT NOT NULL COMMENT '用户id',
    problem_id BIGINT NOT NULL comment '题目id',
    level tinyint DEFAULT 0 comment '题目难度 1:简单 2:中等 3:困难',

    create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY(id),
    NIQUE INDEX uk_user_problem(uid, problem_id)
)engine = InnoDB charset = utf8mb4;

-- 用户解题统计表
create table if not exists statistics_YYYY(
    uid BIGINT not null comment '用户id',

    period CHAR(7) NOT NULL COMMENT 'YYYY-MM',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    submit_count INT DEFAULT 0 comment '题目提交数量',
    accomplish_count INT DEFAULT 0 comment '题目通过数量',
    easy_problem_count INT DEFAULT 0 comment '通过的简单题目数量',
    medium_problem_count INT DEFAULT 0 comment '通过的中等题目数量',
    hard_problem_count INT DEFAULT 0 comment '通过的困难题目数量',

    PRIMARY KEY(period, uid), -- 复合主键
    INDEX idx_accomplish_sort (period, accomplish_count DESC, uid)
)engine = InnoDB charset = utf8mb4;

-- 统计表分区
ALTER TABLE statistics_YYYY
PARTITION BY RANGE COLUMNS(period) (
    PARTITION p202501 VALUES LESS THAN ('2025-02'),
    PARTITION p202502 VALUES LESS THAN ('2025-03'),
    PARTITION p202503 VALUES LESS THAN ('2025-04'),
    PARTITION p202504 VALUES LESS THAN ('2025-05'),
    PARTITION p202505 VALUES LESS THAN ('2025-06'),
    PARTITION p202506 VALUES LESS THAN ('2025-07'),
    PARTITION p202507 VALUES LESS THAN ('2025-08'),
    PARTITION p202508 VALUES LESS THAN ('2025-09'),
    PARTITION p202509 VALUES LESS THAN ('2025-10'),
    PARTITION p202510 VALUES LESS THAN ('2025-11'),
    PARTITION p202511 VALUES LESS THAN ('2025-12'),
    PARTITION p202512 VALUES LESS THAN ('2026-01'),
    PARTITION p_future VALUES LESS THAN (MAXVALUE)
);
