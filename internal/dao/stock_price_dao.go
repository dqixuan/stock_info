package dao

import (
	"context"
	"errors"

	"github.com/dqixuan/stock_info/internal/data"
	"github.com/dqixuan/stock_info/internal/model"
	"gorm.io/gorm"
)

// StockPriceDao 股票价格数据访问对象
type StockPriceDao struct {
	db *gorm.DB
}

// NewStockPriceDao 创建 StockPriceDao
func NewStockPriceDao(data *data.Data) *StockPriceDao {
	return &StockPriceDao{
		db: data.DB(),
	}
}

// Create 创建股票价格记录
func (d *StockPriceDao) Create(ctx context.Context, price *model.StockPrice) error {
	return d.db.WithContext(ctx).Create(price).Error
}

// GetByID 根据ID获取股票价格
func (d *StockPriceDao) GetByID(ctx context.Context, id uint64) (*model.StockPrice, error) {
	var price model.StockPrice
	err := d.db.WithContext(ctx).Where("ID = ?", id).First(&price).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &price, nil
}

// GetByStockIDAndDate 根据股票代码和交易日期获取价格
func (d *StockPriceDao) GetByStockIDAndDate(ctx context.Context, stockID, tradeDate string) (*model.StockPrice, error) {
	var price model.StockPrice
	err := d.db.WithContext(ctx).Where("STOCK_ID = ? AND TRADE_DATE = ?", stockID, tradeDate).First(&price).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &price, nil
}

// ListByStockID 根据股票代码获取价格历史，支持分页
func (d *StockPriceDao) ListByStockID(ctx context.Context, stockID string, page, pageSize int) ([]*model.StockPrice, int64, error) {
	var prices []*model.StockPrice
	var total int64

	query := d.db.WithContext(ctx).Model(&model.StockPrice{}).Where("STOCK_ID = ?", stockID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("TRADE_DATE DESC").Find(&prices).Error; err != nil {
		return nil, 0, err
	}
	return prices, total, nil
}

// ListByStockIDAndDateRange 根据股票代码和时间范围获取价格
func (d *StockPriceDao) ListByStockIDAndDateRange(ctx context.Context, stockID, startDate, endDate string) ([]*model.StockPrice, error) {
	var prices []*model.StockPrice
	err := d.db.WithContext(ctx).Where("STOCK_ID = ? AND TRADE_DATE BETWEEN ? AND ?", stockID, startDate, endDate).
		Order("TRADE_DATE ASC").
		Find(&prices).Error
	return prices, err
}

// ListByTradeDate 根据交易日期获取所有股票价格
func (d *StockPriceDao) ListByTradeDate(ctx context.Context, tradeDate string) ([]*model.StockPrice, error) {
	var prices []*model.StockPrice
	err := d.db.WithContext(ctx).Where("TRADE_DATE = ?", tradeDate).Find(&prices).Error
	return prices, err
}

// List 获取股票价格列表，支持分页
func (d *StockPriceDao) List(ctx context.Context, page, pageSize int) ([]*model.StockPrice, int64, error) {
	var prices []*model.StockPrice
	var total int64

	query := d.db.WithContext(ctx).Model(&model.StockPrice{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("ID ASC").Find(&prices).Error; err != nil {
		return nil, 0, err
	}
	return prices, total, nil
}

// Update 更新股票价格记录
func (d *StockPriceDao) Update(ctx context.Context, price *model.StockPrice) error {
	return d.db.WithContext(ctx).Save(price).Error
}

// UpdateFields 更新指定字段
func (d *StockPriceDao) UpdateFields(ctx context.Context, id uint64, fields map[string]interface{}) error {
	return d.db.WithContext(ctx).Model(&model.StockPrice{}).Where("ID = ?", id).Updates(fields).Error
}

// Delete 删除股票价格记录
func (d *StockPriceDao) Delete(ctx context.Context, id uint64) error {
	return d.db.WithContext(ctx).Delete(&model.StockPrice{}, id).Error
}

// DeleteByStockID 删除某只股票的所有价格记录
func (d *StockPriceDao) DeleteByStockID(ctx context.Context, stockID string) error {
	return d.db.WithContext(ctx).Where("STOCK_ID = ?", stockID).Delete(&model.StockPrice{}).Error
}

// BatchCreate 批量创建股票价格记录
func (d *StockPriceDao) BatchCreate(ctx context.Context, prices []*model.StockPrice) error {
	return d.db.WithContext(ctx).CreateInBatches(prices, 100).Error
}

// Upsert 插入或更新股票价格（根据唯一索引 uk_stock_date）
func (d *StockPriceDao) Upsert(ctx context.Context, price *model.StockPrice) error {
	return d.db.WithContext(ctx).Where("STOCK_ID = ? AND TRADE_DATE = ?", price.StockID, price.TradeDate).
		Assign(price).
		FirstOrCreate(price).Error
}

// BatchUpsert 批量插入或更新股票价格
func (d *StockPriceDao) BatchUpsert(ctx context.Context, prices []*model.StockPrice) error {
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, price := range prices {
			if err := tx.Where("STOCK_ID = ? AND TRADE_DATE = ?", price.StockID, price.TradeDate).
				Assign(price).
				FirstOrCreate(price).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// GetLatestByStockID 获取某只股票最新的价格记录
func (d *StockPriceDao) GetLatestByStockID(ctx context.Context, stockID string) (*model.StockPrice, error) {
	var price model.StockPrice
	err := d.db.WithContext(ctx).Where("STOCK_ID = ?", stockID).
		Order("TRADE_DATE DESC").
		First(&price).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &price, nil
}