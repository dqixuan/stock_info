package model

import "time"

// Stock 股票信息表
type Stock struct {
	ID        uint64    `gorm:"column:ID;primaryKey;autoIncrement" json:"id"`
	StockID   string    `gorm:"column:STOCK_ID;size:20;not null;default:'';uniqueIndex;comment:股票代码" json:"stock_id"`
	Name      string    `gorm:"column:NAME;size:100;not null;default:'';comment:股票名称" json:"name"`
	Market    string    `gorm:"column:MARKET;size:20;not null;default:'';comment:市场(SH/SZ/BJ)" json:"market"`
	Status    int8      `gorm:"column:STATUS;default:1;comment:状态: 1-正常, 2ST  3*ST" json:"status"`
	CreatedAt time.Time `gorm:"column:CREATED_AT;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:UPDATED_AT;autoUpdateTime" json:"updated_at"`
}

// TableName 指定表名
func (Stock) TableName() string {
	return "stocks"
}

// StockPrice 股票价格历史表
type StockPrice struct {
	ID             uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	StockID        uint64    `gorm:"column:stock_id;not null;comment:股票ID" json:"stock_id"`
	Symbol         string    `gorm:"column:symbol;size:20;not null;comment:股票代码" json:"symbol"`
	TradeDate      string    `gorm:"column:trade_date;type:date;not null;comment:交易日期" json:"trade_date"`
	OpenPrice      float64   `gorm:"column:open_price;type:decimal(10,2);not null;comment:开盘价" json:"open_price"`
	ClosePrice     float64   `gorm:"column:close_price;type:decimal(10,2);not null;comment:收盘价" json:"close_price"`
	HighPrice      float64   `gorm:"column:high_price;type:decimal(10,2);not null;comment:最高价" json:"high_price"`
	LowPrice       float64   `gorm:"column:low_price;type:decimal(10,2);not null;comment:最低价" json:"low_price"`
	Volume         int64     `gorm:"column:volume;default:0;comment:成交量" json:"volume"`
	Amount         float64   `gorm:"column:amount;type:decimal(20,2);default:0.00;comment:成交额" json:"amount"`
	ChangePercent  float64   `gorm:"column:change_percent;type:decimal(8,4);default:0.00;comment:涨跌幅(%)" json:"change_percent"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

// TableName 指定表名
func (StockPrice) TableName() string {
	return "stock_prices"
}