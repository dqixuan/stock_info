package model

import "time"

// IndustryStock 行业股票关联表
type IndustryStock struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement;column:ID" json:"id"`
	IndustryID uint64    `gorm:"column:INDUSTRY_ID;type:bigint(20);not null" json:"industry_id"`
	StockID    string    `gorm:"column:STOCK_ID;type:varchar(20);not null;default:''" json:"stock_id"`
	CreatedAt  time.Time `gorm:"column:CREATED_AT;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:UPDATED_AT;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName 指定表名
func (IndustryStock) TableName() string {
	return "industry_stock"
}
