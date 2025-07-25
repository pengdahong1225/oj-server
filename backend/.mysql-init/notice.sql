use oj_online_server;
create table if not exists notice(
    id BIGINT AUTO_INCREMENT,
    create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    delete_at TIMESTAMP,

    title VARCHAR(64) NOT null comment '告示标题',
    content TEXT NOT null comment '告示内容',
    create_by BIGINT DEFAULT 0 comment '题目创建者',
    status tinyint DEFAULT 1 comment '状态 1：正常 0：隐藏',

    PRIMARY KEY(id),
    INDEX idx_title(title)
)engine = InnoDB charset = utf8mb4;