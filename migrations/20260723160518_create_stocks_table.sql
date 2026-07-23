-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS stocks (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    symbol VARCHAR(20) NOT NULL COMMENT '股票代码',
    name VARCHAR(100) NOT NULL COMMENT '股票名称',
    market VARCHAR(20) NOT NULL COMMENT '市场(SH/SZ/HK/US)',
    current_price DECIMAL(10, 2) DEFAULT 0.00 COMMENT '当前价格',
    open_price DECIMAL(10, 2) DEFAULT 0.00 COMMENT '开盘价',
    high_price DECIMAL(10, 2) DEFAULT 0.00 COMMENT '最高价',
    low_price DECIMAL(10, 2) DEFAULT 0.00 COMMENT '最低价',
    volume BIGINT DEFAULT 0 COMMENT '成交量',
    market_cap BIGINT DEFAULT 0 COMMENT '市值',
    status TINYINT DEFAULT 1 COMMENT '状态: 1-正常, 0-停牌',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_symbol_market (symbol, market),
    KEY idx_symbol (symbol),
    KEY idx_market (market)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='股票信息表';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS stocks;
-- +goose StatementEnd
