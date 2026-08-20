-- +goose Up
-- +goose StatementBegin
CREATE TABLE `industry` (
  `ID` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `INDUSTRY_ID` varchar(255) NOT NULL DEFAULT '' COMMENT '行业代码',
  `Name` varchar(255) NOT NULL DEFAULT '' COMMENT '行业名称',
  `COUNT` int(11) NOT NULL DEFAULT 0 COMMENT '行业股票数量',
  `CREATED_AT` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `UPDATED_AT` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`ID`),
  UNIQUE KEY `uk_industry_code` (`INDUSTRY_ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='行业信息表';


CREATE TABLE `industry_stock` (
  `ID` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `INDUSTRY_ID` bigint(20) unsigned NOT NULL COMMENT '行业ID',
  `STOCK_ID` varchar(20) NOT NULL DEFAULT '' COMMENT '股票代码',
  `CREATED_AT` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `UPDATED_AT` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`ID`),
  UNIQUE KEY `uk_industry_stock` (`INDUSTRY_ID`, `STOCK_ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='行业股票关联表';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
