package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	v1 "github.com/dqixuan/stock_info/api/stock/v1"
	"github.com/dqixuan/stock_info/internal/dao"
	"github.com/dqixuan/stock_info/internal/data"
	"github.com/dqixuan/stock_info/internal/model"
	"github.com/dqixuan/stock_info/internal/utils"
	"github.com/dqixuan/stock_info/pkg"
	"gorm.io/gorm"
)

type StockService struct {
	mysql         *gorm.DB
	stockDao      *dao.StockDao
	stockPriceDao *dao.StockPriceDao
}

func (s StockService) InitStockInfo(ctx context.Context, request *v1.InitStockInfoRequest) (*v1.InitStockInfoReply, error) {
	fn := "InitStockInfo"
	go func() {
		stocks, err := pkg.FetchAllAStocks()
		if err != nil {
			fmt.Println(fn, "err:", err)
			return
		}
		ctxWithoutCancel := context.WithoutCancel(ctx)
		status := model.StockStatusNormal
		for _, stock := range stocks {
			if strings.Contains(stock.Name, "ST") {
				status = model.StockStatusST
			}
			if strings.Contains(stock.Name, "*ST") {
				status = model.StockStatusStarST
			}
			err = s.stockDao.Create(ctxWithoutCancel, &model.Stock{
				StockID:   stock.Code,
				Name:      stock.Name,
				ShortName: utils.GetInitials(stock.Name),
				Market:    stock.Symbol[:2],
				Status:    int8(status),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			})
			if err != nil {
				fmt.Println(fn, "err:", err)
			}
		}

	}()
	return &v1.InitStockInfoReply{}, nil
}

func (s StockService) AddStockPrice(ctx context.Context, request *v1.StockPriceRequest) (*v1.StockPriceReply, error) {
	price := &model.StockPrice{
		StockID:       request.StockId,
		TradeDate:     request.TradeDate,
		OpenPrice:     request.OpenPrice,
		ClosePrice:    request.ClosePrice,
		HighPrice:     request.HighPrice,
		LowPrice:      request.LowPrice,
		Volume:        request.Volume,
		Amount:        request.Amount,
		ChangePercent: request.ChangePercent,
		MarginBalance: request.MarginBalance,
	}

	if err := s.stockPriceDao.Upsert(ctx, price); err != nil {
		return nil, err
	}

	return &v1.StockPriceReply{}, nil
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
