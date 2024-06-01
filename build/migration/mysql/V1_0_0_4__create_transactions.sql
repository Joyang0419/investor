CREATE TABLE IF NOT EXISTS transactions
(
    transaction_id INT AUTO_INCREMENT PRIMARY KEY COMMENT '交易ID',
    type ENUM('deposit', 'withdrawal', 'transfer') NOT NULL COMMENT '交易类型',
    amount DECIMAL(10, 2) NOT NULL COMMENT '交易金额',
    account_id INT NOT NULL COMMENT '交易账户ID',
    target_account_id INT COMMENT '对于转账，需要目标账户ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '交易时间'
);
