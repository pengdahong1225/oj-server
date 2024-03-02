use gvb_server;

-- 用户信息表
create table if not exists user_data
(
    id BIGINT AUTO_INCREMENT,
    create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    delete_at TIMESTAMP,

    phone BIGINT not NULL,
    pwd VARCHAR(64) NOT null,
    nickname VARCHAR(64) DEFAULT '新用户',
    email VARCHAR(64) DEFAULT '',
    gender tinyint DEFAULT 0 comment '0:woman 1:man',
    role tinyint DEFAULT 0 comment '0:user 1:admin',
    head_pic VARCHAR(256) DEFAULT '' comment '头像url',

    PRIMARY KEY(id),
    UNIQUE INDEX idx_phone(phone)
)engine = InnoDB charset = utf8mb4;

-- 文章类别表
create table if not exists article_cate
(
    cate_id BIGINT AUTO_INCREMENT,
    description VARCHAR(64) NOT null,

    PRIMARY KEY(cate_id)
)engine = InnoDB charset = utf8mb4;

-- 文章表
create table if not exists article
(
    id BIGINT AUTO_INCREMENT,
    title VARCHAR(64) NOT null,
    state tinyint DEFAULT 1,    -- 0:关 1:开
    cate_id BIGINT,
    text TEXT,      -- 正文
    abstract TEXT,  -- 摘要
    create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY(id),
    FOREIGN KEY (cate_id) REFERENCES article_cate(cate_id)
)engine = InnoDB charset = utf8mb4;

