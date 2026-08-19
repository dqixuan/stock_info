package model

import "time"

// StockPrice 股票价格历史表
type StockPrice struct {
	ID                 uint64    `gorm:"column:ID;primaryKey;autoIncrement" json:"id"`
	StockID            string    `gorm:"column:STOCK_ID;size:20;not null;default:'';comment:股票代码" json:"stock_id"`
	TradeDate          string    `gorm:"column:TRADE_DATE;type:date;not null;comment:交易日期" json:"trade_date"`
	OpenPrice          float64   `gorm:"column:OPEN_PRICE;type:decimal(10,2);not null;default:0.00;comment:开盘价" json:"open_price"`
	ClosePrice         float64   `gorm:"column:CLOSE_PRICE;type:decimal(10,2);not null;default:0.00;comment:收盘价" json:"close_price"`
	HighPrice          float64   `gorm:"column:HIGH_PRICE;type:decimal(10,2);not null;default:0.00;comment:最高价" json:"high_price"`
	LowPrice           float64   `gorm:"column:LOW_PRICE;type:decimal(10,2);not null;default:0.00;comment:最低价" json:"low_price"`
	Volume             int64     `gorm:"column:VOLUME;default:0;comment:成交量" json:"volume"`
	Amount             float64   `gorm:"column:AMOUNT;type:decimal(20,2);default:0.00;comment:成交额" json:"amount"`
	ChangePercent      float64   `gorm:"column:CHANGE_PERCENT;type:decimal(8,4);default:0.00;comment:涨跌幅(%)" json:"change_percent"`
	FinanceBalance     float64   `gorm:"column:FINANCE_BALANCE;type:decimal(20,2);default:0.00;comment:融资余额" json:"finance_balance"`
	FinanceBuy         float64   `gorm:"column:FINANCE_BUY;type:decimal(20,2);default:0.00;comment:融资买入额" json:"finance_buy"`
	FinanceRepay       float64   `gorm:"column:FINANCE_REPAY;type:decimal(20,2);default:0.00;comment:融资偿还额" json:"finance_repay"`
	SecurityLendVolume float64   `gorm:"column:SECURITY_LEND_VOLUME;type:decimal(20,2);default:0.00;comment:融券余量" json:"security_lend_volume"`
	SecurityLendSell   float64   `gorm:"column:SECURITY_LEND_SELL;type:decimal(20,2);default:0.00;comment:融券卖出量" json:"security_lend_sell"`
	SecurityLendRepay  float64   `gorm:"column:SECURITY_LEND_REPAY;type:decimal(20,2);default:0.00;comment:融券偿还量" json:"security_lend_repay"`
	CreatedAt          time.Time `gorm:"column:CREATED_AT;autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time `gorm:"column:UPDATED_AT;autoUpdateTime" json:"updated_at"`
}

// TableName 指定表名
func (StockPrice) TableName() string {
	return "stock_price"
}
