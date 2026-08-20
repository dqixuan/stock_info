package dao

import (
	"context"

	"github.com/dqixuan/stock_info/internal/data"
	"github.com/dqixuan/stock_info/internal/model"
	"gorm.io/gorm"
)

// ==========================================
// 2. IndustryStock DAO (行业股票关联表)
// ==========================================

type IndustryStockDao struct {
	db *gorm.DB
}

func NewIndustryStockDao(data *data.Data) *IndustryStockDao {
	return &IndustryStockDao{db: data.DB()}
}

// BatchCreate 批量添加行业与股票的关联
func (d *IndustryStockDao) BatchCreate(ctx context.Context, stocks []model.IndustryStock) error {
	if len(stocks) == 0 {
		return nil
	}
	// 关联表数据量可能较大，每批次 500 条
	return d.db.WithContext(ctx).CreateInBatches(stocks, 500).Error
}

// DeleteByIndustryID 删除某行业下的所有股票关联 (常用于行业成分股全量更新前的清理)
func (d *IndustryStockDao) DeleteByIndustryID(ctx context.Context, industryID uint64) error {
	return d.db.WithContext(ctx).Where("INDUSTRY_ID = ?", industryID).Delete(&model.IndustryStock{}).Error
}

// DeleteByStockID 删除某股票的所有行业关联 (常用于退市股票清理)
func (d *IndustryStockDao) DeleteByStockID(ctx context.Context, stockID string) error {
	return d.db.WithContext(ctx).Where("STOCK_ID = ?", stockID).Delete(&model.IndustryStock{}).Error
}

// GetStocksByIndustryID 查询某行业下的所有股票 (正向查询)
func (d *IndustryStockDao) GetStocksByIndustryID(ctx context.Context, industryID uint64) ([]model.IndustryStock, error) {
	var stocks []model.IndustryStock
	err := d.db.WithContext(ctx).Where("INDUSTRY_ID = ?", industryID).Find(&stocks).Error
	return stocks, err
}

// GetStockCodesByIndustryID 查询某行业下的所有股票代码 (只提取代码，节省内存)
func (d *IndustryStockDao) GetStockCodesByIndustryID(ctx context.Context, industryID uint64) ([]string, error) {
	var stockCodes []string
	err := d.db.WithContext(ctx).Model(&model.IndustryStock{}).
		Where("INDUSTRY_ID = ?", industryID).
		Pluck("STOCK_ID", &stockCodes).Error
	return stockCodes, err
}

// GetIndustryIDsByStockID 查询某股票所属的所有行业ID (反向查询)
func (d *IndustryStockDao) GetIndustryIDsByStockID(ctx context.Context, stockID string) ([]uint64, error) {
	var industryIDs []uint64
	err := d.db.WithContext(ctx).Model(&model.IndustryStock{}).
		Where("STOCK_ID = ?", stockID).
		Pluck("INDUSTRY_ID", &industryIDs).Error
	return industryIDs, err
}
