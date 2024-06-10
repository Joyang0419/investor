ALTER TABLE transactions
    MODIFY COLUMN type ENUM('deposit', 'withdraw', 'transfer') NOT NULL COMMENT '交易类型';
