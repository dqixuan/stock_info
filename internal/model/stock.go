package model

import "time"

// Stock 股票信息表
type Stock struct {
	ID        uint64    `gorm:"column:ID;primaryKey;autoIncrement" json:"id"`
	StockID   string    `gorm:"column:STOCK_ID;size:20;not null;default:'';uniqueIndex;comment:股票代码" json:"stock_id"`
	Name      string    `gorm:"column:NAME;size:100;not null;default:'';comment:股票名称" json:"name"`
	ShortName string    `gorm:"column:SHORT_NAME;size:50;not null;default:'';comment:股票简称" json:"short_name"`
	Market    string    `gorm:"column:MARKET;size:20;not null;default:'';comment:市场(SH/SZ/BJ)" json:"market"`
	Status    int8      `gorm:"column:STATUS;default:1;comment:状态: 1-正常, 2ST  3*ST" json:"status"`
	CreatedAt time.Time `gorm:"column:CREATED_AT;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:UPDATED_AT;autoUpdateTime" json:"updated_at"`
}

// TableName 指定表名
func (Stock) TableName() string {
	return "stocks"
}
