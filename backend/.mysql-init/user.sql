use oj_online_server;

create table if not exists user_info
(
    id BIGINT AUTO_INCREMENT,
    create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    delete_at TIMESTAMP,

    mobile BIGINT not NULL comment '手机号',
    password VARCHAR(256) not NULL comment '密码',
    nickname VARCHAR(64) DEFAULT '新用户',
    email VARCHAR(64) DEFAULT '',
    gender tinyint DEFAULT 0 comment '0:woman 1:man',
    role tinyint DEFAULT 0 comment '0:user 1:admin',
    avatar_url VARCHAR(256) DEFAULT '' comment '头像url',

    PRIMARY KEY(id),
    UNIQUE INDEX idx_mobile(mobile)
)engine = InnoDB charset = utf8mb4;

INSERT IGNORE INTO user_info (mobile, password, nickname, email, gender, role, avatar_url)
VALUES (13800000000, 'admin123', 'admin', 'admin@example.com', 1, 1, '');