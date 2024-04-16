-- 创建表格
CREATE TABLE IF NOT EXISTS dailyPrice
(
    id             INT AUTO_INCREMENT PRIMARY KEY,
    stock_code     VARCHAR(255) NOT NULL,
    highest_price  FLOAT        NOT NULL,
    lowest_price   FLOAT        NOT NULL,
    opening_price  FLOAT        NOT NULL,
    closing_price  FLOAT        NOT NULL,
    volume         BIGINT       NOT NULL,
    `change`       FLOAT        NOT NULL,
    date_timestamp BIGINT       NOT NULL
);