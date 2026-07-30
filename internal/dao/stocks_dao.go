package dao

import (
	"errors"

	"github.com/dqixuan/stock_info/internal/model"
	"gorm.io/gorm"
)

// StockDao 股票信息数据访问对象
type StockDao struct {
	db *gorm.DB
}

// NewStockDao 创建 StockDao
func NewStockDao(db *gorm.DB) *StockDao {
	return &StockDao{db: db}
}

// Create 创建股票
func (d *StockDao) Create(stock *model.Stock) error {
	return d.db.Create(stock).Error
}

// GetByID 根据ID获取股票
func (d *StockDao) GetByID(id uint64) (*model.Stock, error) {
	var stock model.Stock
	err := d.db.Where("ID = ?", id).First(&stock).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &stock, nil
}

// GetByStockID 根据股票代码获取股票
func (d *StockDao) GetByStockID(stockID string) (*model.Stock, error) {
	var stock model.Stock
	err := d.db.Where("STOCK_ID = ?", stockID).First(&stock).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &stock, nil
}

// List 获取股票列表，支持分页
func (d *StockDao) List(page, pageSize int) ([]*model.Stock, int64, error) {
	var stocks []*model.Stock
	var total int64

	query := d.db.Model(&model.Stock{})
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
func (d *StockDao) ListByMarket(market string, page, pageSize int) ([]*model.Stock, int64, error) {
	var stocks []*model.Stock
	var total int64

	query := d.db.Model(&model.Stock{}).Where("MARKET = ?", market)
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
func (d *StockDao) Update(stock *model.Stock) error {
	return d.db.Save(stock).Error
}

// UpdateFields 更新指定字段
func (d *StockDao) UpdateFields(id uint64, fields map[string]interface{}) error {
	return d.db.Model(&model.Stock{}).Where("ID = ?", id).Updates(fields).Error
}

// Delete 删除股票
func (d *StockDao) Delete(id uint64) error {
	return d.db.Delete(&model.Stock{}, id).Error
}

// BatchCreate 批量创建股票
func (d *StockDao) BatchCreate(stocks []*model.Stock) error {
	return d.db.CreateInBatches(stocks, 100).Error
}

// Upsert 插入或更新股票（根据唯一索引 STOCK_ID）
func (d *StockDao) Upsert(stock *model.Stock) error {
	return d.db.Where("STOCK_ID = ?", stock.StockID).Assign(stock).FirstOrCreate(stock).Error
}

// StockPriceDao 股票价格数据访问对象
type StockPriceDao struct {
	db *gorm.DB
}

// NewStockPriceDao 创建 StockPriceDao
func NewStockPriceDao(db *gorm.DB) *StockPriceDao {
	return &StockPriceDao{db: db}
}

// Create 创建股票价格记录
func (d *StockPriceDao) Create(price *model.StockPrice) error {
	return d.db.Create(price).Error
}

// GetByID 根据ID获取股票价格
func (d *StockPriceDao) GetByID(id uint64) (*model.StockPrice, error) {
	var price model.StockPrice
	err := d.db.First(&price, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &price, nil
}

// GetByStockIDAndDate 根据股票ID和交易日期获取价格
func (d *StockPriceDao) GetByStockIDAndDate(stockID uint64, tradeDate string) (*model.StockPrice, error) {
	var price model.StockPrice
	err := d.db.Where("stock_id = ? AND trade_date = ?", stockID, tradeDate).First(&price).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &price, nil
}

// ListByStockID 根据股票ID获取价格历史
func (d *StockPriceDao) ListByStockID(stockID uint64, page, pageSize int) ([]*model.StockPrice, int64, error) {
	var prices []*model.StockPrice
	var total int64

	query := d.db.Model(&model.StockPrice{}).Where("stock_id = ?", stockID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("trade_date DESC").Find(&prices).Error; err != nil {
		return nil, 0, err
	}
	return prices, total, nil
}

// ListBySymbolAndDateRange 根据股票代码和时间范围获取价格
func (d *StockPriceDao) ListBySymbolAndDateRange(symbol string, startDate, endDate string) ([]*model.StockPrice, error) {
	var prices []*model.StockPrice
	err := d.db.Where("symbol = ? AND trade_date BETWEEN ? AND ?", symbol, startDate, endDate).
		Order("trade_date ASC").
		Find(&prices).Error
	return prices, err
}

// ListByTradeDate 根据交易日期获取所有股票价格
func (d *StockPriceDao) ListByTradeDate(tradeDate string) ([]*model.StockPrice, error) {
	var prices []*model.StockPrice
	err := d.db.Where("trade_date = ?", tradeDate).Find(&prices).Error
	return prices, err
}

// Update 更新股票价格记录
func (d *StockPriceDao) Update(price *model.StockPrice) error {
	return d.db.Save(price).Error
}

// UpdateFields 更新指定字段
func (d *StockPriceDao) UpdateFields(id uint64, fields map[string]interface{}) error {
	return d.db.Model(&model.StockPrice{}).Where("id = ?", id).Updates(fields).Error
}

// Delete 删除股票价格记录
func (d *StockPriceDao) Delete(id uint64) error {
	return d.db.Delete(&model.StockPrice{}, id).Error
}

// DeleteByStockID 删除某只股票的所有价格记录
func (d *StockPriceDao) DeleteByStockID(stockID uint64) error {
	return d.db.Where("stock_id = ?", stockID).Delete(&model.StockPrice{}).Error
}

// BatchCreate 批量创建股票价格记录
func (d *StockPriceDao) BatchCreate(prices []*model.StockPrice) error {
	return d.db.CreateInBatches(prices, 100).Error
}

// Upsert 插入或更新股票价格（根据唯一索引 uk_stock_date）
func (d *StockPriceDao) Upsert(price *model.StockPrice) error {
	return d.db.Where("stock_id = ? AND trade_date = ?", price.StockID, price.TradeDate).
		Assign(price).
		FirstOrCreate(price).Error
}

// BatchUpsert 批量插入或更新股票价格
func (d *StockPriceDao) BatchUpsert(prices []*model.StockPrice) error {
	return d.db.Transaction(func(tx *gorm.DB) error {
		for _, price := range prices {
			if err := tx.Where("stock_id = ? AND trade_date = ?", price.StockID, price.TradeDate).
				Assign(price).
				FirstOrCreate(price).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// GetLatestByStockID 获取某只股票最新的价格记录
func (d *StockPriceDao) GetLatestByStockID(stockID uint64) (*model.StockPrice, error) {
	var price model.StockPrice
	err := d.db.Where("stock_id = ?", stockID).
		Order("trade_date DESC").
		First(&price).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &price, nil
}