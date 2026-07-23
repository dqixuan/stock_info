-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS stock_prices (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    stock_id BIGINT NOT NULL COMMENT '股票ID',
    symbol VARCHAR(20) NOT NULL COMMENT '股票代码',
    trade_date DATE NOT NULL COMMENT '交易日期',
    open_price DECIMAL(10, 2) NOT NULL COMMENT '开盘价',
    close_price DECIMAL(10, 2) NOT NULL COMMENT '收盘价',
    high_price DECIMAL(10, 2) NOT NULL COMMENT '最高价',
    low_price DECIMAL(10, 2) NOT NULL COMMENT '最低价',
    volume BIGINT DEFAULT 0 COMMENT '成交量',
    amount DECIMAL(20, 2) DEFAULT 0.00 COMMENT '成交额',
    change_percent DECIMAL(8, 4) DEFAULT 0.00 COMMENT '涨跌幅(%)',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uk_stock_date (stock_id, trade_date),
    KEY idx_symbol_date (symbol, trade_date),
    KEY idx_trade_date (trade_date),
    FOREIGN KEY (stock_id) REFERENCES stocks(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='股票价格历史表';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS stock_prices;
-- +goose StatementEnd
