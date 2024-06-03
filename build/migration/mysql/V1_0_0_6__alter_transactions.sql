ALTER TABLE transactions
    ADD COLUMN completed BOOLEAN NOT NULL DEFAULT FALSE COMMENT '交易完成' AFTER target_account_id;
