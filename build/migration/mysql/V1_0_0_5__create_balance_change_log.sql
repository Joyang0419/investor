CREATE TABLE IF NOT EXISTS balance_change_log
(
    account_id     INT            NOT NULL COMMENT '交易帳戶ID',
    order_id       INT            NOT NULL COMMENT '订单ID',
    before_balance DECIMAL(10, 2) NOT NULL COMMENT '交易前余额',
    after_balance  DECIMAL(10, 2) NOT NULL COMMENT '交易后余额',
    created_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT 'log產生时间',
    PRIMARY KEY (account_id, order_id) COMMENT '主键'
);
