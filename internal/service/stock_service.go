package service

import (
	"context"

	v1 "github.com/dqixuan/stock_info/api/stock/v1"
	"github.com/dqixuan/stock_info/internal/dao"
	"github.com/dqixuan/stock_info/internal/data"
	"gorm.io/gorm"
)

type StockService struct {
	mysql         *gorm.DB
	stockDao      *dao.StockDao
	stockPriceDao *dao.StockPriceDao
}

func (s StockService) AddStockPrice(ctx context.Context, request *v1.StockPriceRequest) (*v1.StockPriceReply, error) {
	//TODO implement me
	panic("implement me")
}

func (s StockService) GetStockList(ctx context.Context, request *v1.StockInfoRequest) (*v1.StockInfoReply, error) {
	//TODO implement me
	panic("implement me")
}

func NewStockService(
	data *data.Data,
	stockDao *dao.StockDao,
	stockPriceDao *dao.StockPriceDao,
) *StockService {
	return &StockService{
		mysql:         data.DB(),
		stockDao:      stockDao,
		stockPriceDao: stockPriceDao,
	}
}

var _ v1.StockServiceHTTPServer = (*StockService)(nil)
