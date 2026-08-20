package model

import "time"

// Industry 行业信息表
type Industry struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement;column:ID" json:"id"`
	IndustryID string    `gorm:"column:INDUSTRY_ID;uniqueIndex:uk_industry_code;type:varchar(255);not null;default:''" json:"industry_id"`
	Name       string    `gorm:"column:Name;type:varchar(255);not null;default:''" json:"name"`
	Count      int       `gorm:"column:COUNT;type:int(11);not null;default:0" json:"count"`
	CreatedAt  time.Time `gorm:"column:CREATED_AT;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:UPDATED_AT;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName 指定表名，防止 GORM 自动将其转换为复数形式 (industries)
func (Industry) TableName() string {
	return "industry"
}
