CREATE DATABASE IF NOT EXISTS `open_oj` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;

USE `open_oj`;

-- 用户信息表
CREATE TABLE IF NOT EXISTS user_info
(
    id         BIGINT AUTO_INCREMENT,
    create_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
    update_at  DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    delete_at  DATETIME NULL DEFAULT NULL COMMENT '软删除时间',

    mobile     BIGINT       NOT NULL COMMENT '手机号',
    password   VARCHAR(256) NOT NULL COMMENT '密码(哈希)',
    username   VARCHAR(64)  DEFAULT '新用户' COMMENT '用户名',
    email      VARCHAR(64)  DEFAULT '',
    role       TINYINT      DEFAULT 0 COMMENT '0:user 1:admin',
    avatar_url VARCHAR(256) DEFAULT '' COMMENT '头像url',

    PRIMARY KEY (id),
    UNIQUE INDEX idx_mobile (mobile),
    INDEX idx_delete_at (delete_at)
) ENGINE = InnoDB CHARSET = utf8mb4;

INSERT IGNORE INTO user_info (mobile, password, username, email, role, avatar_url)
VALUES (18000000000, 'admin123', 'admin', 'admin@example.com', 1, '');

-- 题目表
CREATE TABLE IF NOT EXISTS problem
(
    id            BIGINT AUTO_INCREMENT,
    create_at     DATETIME DEFAULT CURRENT_TIMESTAMP,
    update_at     DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    delete_at     DATETIME NULL DEFAULT NULL COMMENT '软删除时间',

    title         VARCHAR(64)  NOT NULL COMMENT '题目标题',
    level         TINYINT      DEFAULT 0 COMMENT '题目难度 1:简单 2:中等 3:困难',
    tags          JSON         COMMENT '题目标签',
    package_dir   VARCHAR(256) DEFAULT '' COMMENT '题目包目录',
    status        TINYINT      DEFAULT 0 COMMENT '状态 0：待校验 1：校验中 2：未通过校验 3：待发布 4：已发布 5：已下线',

    submit_count  INT DEFAULT 0 COMMENT '提交数',
    accept_count  INT DEFAULT 0 COMMENT '通过数',
    comment_count INT DEFAULT 0 COMMENT '评论数',

    PRIMARY KEY (id),
    UNIQUE INDEX idx_title (title),
    INDEX idx_delete_at (delete_at),
    INDEX idx_status (status)
) ENGINE = InnoDB CHARSET = utf8mb4;

-- 告示表
CREATE TABLE IF NOT EXISTS notice
(
    id         BIGINT AUTO_INCREMENT,
    create_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
    update_at  DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    delete_at  DATETIME NULL DEFAULT NULL COMMENT '软删除时间',

    title      VARCHAR(64) NOT NULL COMMENT '告示标题',
    content    TEXT        NOT NULL COMMENT '告示内容',
    create_by  BIGINT      DEFAULT 0 COMMENT '公告创建者',
    status     TINYINT     DEFAULT 0 COMMENT '状态 1：正常 0：隐藏',

    PRIMARY KEY (id),
    INDEX idx_status_create (status, create_at),
    INDEX idx_delete_at (delete_at)
) ENGINE = InnoDB CHARSET = utf8mb4;

-- 评论表
CREATE TABLE IF NOT EXISTS comment
(
    id               BIGINT AUTO_INCREMENT,
    create_at        DATETIME DEFAULT CURRENT_TIMESTAMP,
    update_at        DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    delete_at        DATETIME NULL DEFAULT NULL COMMENT '软删除时间',

    obj_id           BIGINT       NOT NULL COMMENT '对象id（题目id）',
    user_id          BIGINT       NOT NULL COMMENT '评论用户id',
    user_name        VARCHAR(64)  DEFAULT '' COMMENT '评论时用户名',
    user_avatar_url  VARCHAR(256) DEFAULT '' COMMENT '评论时头像url',
    content          TEXT         NOT NULL COMMENT '评论内容',
    status           TINYINT      DEFAULT 1 COMMENT '状态 1：正常 0：隐藏',
    reply_count      INT          DEFAULT 0 COMMENT '回复数量',
    like_count       INT          DEFAULT 0 COMMENT '点赞数量',
    child_count      INT          DEFAULT 0 COMMENT '两层结构，root才有子评论',
    pub_stamp        BIGINT       NOT NULL COMMENT '时间戳',
    pub_region       VARCHAR(64)  DEFAULT '' COMMENT '发布地区',

    is_root          TINYINT      DEFAULT 1 COMMENT '是否是楼主 1：是 0：不是',
    root_id          BIGINT       NOT NULL COMMENT '楼主id',
    root_comment_id  BIGINT       NOT NULL COMMENT '楼主评论id',
    reply_id         BIGINT       NOT NULL COMMENT '回复用户id',
    reply_comment_id BIGINT       NOT NULL COMMENT '回复评论id',
    reply_user_name  VARCHAR(64)  DEFAULT '' COMMENT '回复用户名',

    PRIMARY KEY (id),
    INDEX idx_obj_id (obj_id),
    INDEX idx_user_id (user_id),
    INDEX idx_root_id (root_id),
    INDEX idx_reply_id (reply_id),
    INDEX idx_like_count (like_count),
    INDEX idx_delete_at (delete_at),

    -- 拉取某题目下的评论：WHERE obj_id=? AND status=1 ORDER BY like_count DESC
    INDEX idx_obj_status_like (obj_id, status, like_count),
    INDEX idx_obj_status_stamp (obj_id, status, pub_stamp)
) ENGINE = InnoDB CHARSET = utf8mb4;

-- 用户提交记录表
CREATE TABLE IF NOT EXISTS user_submit_record
(
    id           BIGINT AUTO_INCREMENT,
    created_at   DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at   DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    uid          BIGINT       NOT NULL COMMENT '用户id',
    user_name    VARCHAR(64)  NOT NULL COMMENT '用户名',
    problem_id   BIGINT       NOT NULL COMMENT '题目id',
    problem_name VARCHAR(64)  NOT NULL COMMENT '题目名称',

    lang         VARCHAR(64)  DEFAULT '' COMMENT '语言',
    code         TEXT         NOT NULL COMMENT '提交的代码',

    status       TINYINT      DEFAULT 0 COMMENT '0:排队 1:判题中 2:完成',
    accepted     TINYINT(1)   DEFAULT 0 COMMENT '是否通过',
    message      VARCHAR(256) DEFAULT '' COMMENT '测评结果描述',
    result       MEDIUMBLOB   NOT NULL COMMENT '运行结果集',

    PRIMARY KEY (id),
    INDEX idx_uid_created (uid, created_at),
    INDEX idx_problem_created (problem_id, created_at),
    INDEX idx_uid_problem_created (uid, problem_id, created_at)
) ENGINE = InnoDB CHARSET = utf8mb4;

-- 用户解题表
CREATE TABLE IF NOT EXISTS user_solution
(
    id         BIGINT AUTO_INCREMENT,
    uid        BIGINT  NOT NULL COMMENT '用户id',
    problem_id BIGINT  NOT NULL COMMENT '题目id',
    level      TINYINT DEFAULT 0 COMMENT '题目难度 1:简单 2:中等 3:困难',

    create_at  DATETIME DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (id),
    UNIQUE INDEX uk_user_problem (uid, problem_id)
) ENGINE = InnoDB CHARSET = utf8mb4;

-- 用户解题统计表
CREATE TABLE IF NOT EXISTS statistics_2026
(
    uid                  BIGINT  NOT NULL COMMENT '用户id',

    period               CHAR(7) NOT NULL COMMENT 'YYYY-MM',
    updated_at           DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    submit_count         INT DEFAULT 0 COMMENT '题目提交数量',
    accomplish_count     INT DEFAULT 0 COMMENT '题目通过数量',
    easy_problem_count   INT DEFAULT 0 COMMENT '通过的简单题目数量',
    medium_problem_count INT DEFAULT 0 COMMENT '通过的中等题目数量',
    hard_problem_count   INT DEFAULT 0 COMMENT '通过的困难题目数量',

    PRIMARY KEY (period, uid),
    INDEX idx_accomplish_sort (period, accomplish_count DESC, uid)
) ENGINE = InnoDB CHARSET = utf8mb4;

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
