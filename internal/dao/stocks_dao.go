package dao

import (
	"context"
	"errors"

	"github.com/dqixuan/stock_info/internal/data"
	"github.com/dqixuan/stock_info/internal/model"
	"gorm.io/gorm"
)

// StockDao 股票信息数据访问对象
type StockDao struct {
	db *gorm.DB
}

// NewStockDao 创建 StockDao
func NewStockDao(data *data.Data) *StockDao {
	return &StockDao{
		db: data.DB(),
	}
}

// Create 创建股票
func (d *StockDao) Create(ctx context.Context, stock *model.Stock) error {
	return d.db.WithContext(ctx).Create(stock).Error
}

// GetByID 根据ID获取股票
func (d *StockDao) GetByID(ctx context.Context, id uint64) (*model.Stock, error) {
	var stock model.Stock
	err := d.db.WithContext(ctx).Where("ID = ?", id).First(&stock).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &stock, nil
}

// GetByStockID 根据股票代码获取股票
func (d *StockDao) GetByStockID(ctx context.Context, stockID string) (*model.Stock, error) {
	var stock model.Stock
	err := d.db.WithContext(ctx).Where("STOCK_ID = ?", stockID).First(&stock).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &stock, nil
}

// List 获取股票列表，支持分页
func (d *StockDao) List(ctx context.Context, page, pageSize int) ([]*model.Stock, int64, error) {
	var stocks []*model.Stock
	var total int64

	query := d.db.WithContext(ctx).Model(&model.Stock{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("ID ASC").Find(&stocks).Error; err != nil {
		return nil, 0, err
	}
	return stocks, total, nil
}

// ListByMarket 根据市场获取股票列表
func (d *StockDao) ListByMarket(ctx context.Context, market string, page, pageSize int) ([]*model.Stock, int64, error) {
	var stocks []*model.Stock
	var total int64

	query := d.db.WithContext(ctx).Model(&model.Stock{}).Where("MARKET = ?", market)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("ID ASC").Find(&stocks).Error; err != nil {
		return nil, 0, err
	}
	return stocks, total, nil
}

// Update 更新股票信息
func (d *StockDao) Update(ctx context.Context, stock *model.Stock) error {
	return d.db.WithContext(ctx).Save(stock).Error
}

// UpdateFields 更新指定字段
func (d *StockDao) UpdateFields(ctx context.Context, id uint64, fields map[string]interface{}) error {
	return d.db.WithContext(ctx).Model(&model.Stock{}).Where("ID = ?", id).Updates(fields).Error
}

// Delete 删除股票
func (d *StockDao) Delete(ctx context.Context, id uint64) error {
	return d.db.WithContext(ctx).Delete(&model.Stock{}, id).Error
}

// BatchCreate 批量创建股票
func (d *StockDao) BatchCreate(ctx context.Context, stocks []*model.Stock) error {
	return d.db.WithContext(ctx).CreateInBatches(stocks, 100).Error
}

// Upsert 插入或更新股票（根据唯一索引 STOCK_ID）
func (d *StockDao) Upsert(ctx context.Context, stock *model.Stock) error {
	return d.db.WithContext(ctx).Where("STOCK_ID = ?", stock.StockID).Assign(stock).FirstOrCreate(stock).Error
}