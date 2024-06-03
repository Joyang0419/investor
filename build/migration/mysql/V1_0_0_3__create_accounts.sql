CREATE TABLE IF NOT EXISTS accounts
(
    id    INT AUTO_INCREMENT PRIMARY KEY COMMENT '帐户ID',
    email         VARCHAR(255)   NOT NULL COMMENT '电子邮件地址',
    name  VARCHAR(255)   NOT NULL COMMENT '用户的全名',
    picture       VARCHAR(255) COMMENT 'GoogleOauth 用户的头像图片 URL',
    balance       DECIMAL(10, 2) NOT NULL DEFAULT 0.00 COMMENT '账户余额',
    last_login_at TIMESTAMP               DEFAULT CURRENT_TIMESTAMP COMMENT '上次登录时间戳',
    created_at    TIMESTAMP               DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间'
);