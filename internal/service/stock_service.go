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

const (
	pageSize    = 100
	workerCount = 5
)

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

func (s StockService) UpdateStockPrice(ctx context.Context, request *v1.UpdateStockPriceRequest) (*v1.UpdateStockPriceReply, error) {

	fn := "UpdateStockPrice"
	go func() {
		ctxWithoutCancel := context.WithoutCancel(ctx)
		page := 1
		for {
			stocks, err := s.stockDao.List(ctxWithoutCancel, page, pageSize)
			if err != nil {
				fmt.Println(fn, "list stocks err:", err)
				return
			}
			if len(stocks) == 0 {
				return
			}

			prices := make([]*model.StockPrice, 0, len(stocks))

			tradeDate := time.Now().Format(time.DateOnly)
			for _, stock := range stocks {
				if stock == nil || stock.StockID == "" {
					continue
				}

				stockData, err := pkg.GetStockPrice(stock.StockID)
				if err != nil {
					fmt.Println(fn, "get stock data err:", stock.StockID, err)
					continue
				}

				margin, err := pkg.GetStockMarginByDate(stock.StockID, compactTradeDate(tradeDate))
				if err != nil {
					fmt.Println(fn, "get stock margin err:", stock.StockID, err)
					margin = nil
				}

				price := &model.StockPrice{
					StockID:            stock.StockID,
					TradeDate:          tradeDate,
					OpenPrice:          stockData.OpenPrice,
					ClosePrice:         stockData.LatestPrice,
					HighPrice:          stockData.HighPrice,
					LowPrice:           stockData.LowPrice,
					Volume:             int64(stockData.Volume),
					Amount:             stockData.Turnover,
					ChangePercent:      stockData.ChangePercentage,
					FinanceBalance:     marginValue(margin, func(m *pkg.MarginData) float64 { return m.MarginBalance }),
					FinanceBuy:         marginValue(margin, func(m *pkg.MarginData) float64 { return m.MarginBuyAmount }),
					FinanceRepay:       marginValue(margin, func(m *pkg.MarginData) float64 { return m.MarginRepayAmount }),
					SecurityLendVolume: marginValue(margin, func(m *pkg.MarginData) float64 { return m.ShortSellVolume }),
					SecurityLendSell:   marginValue(margin, func(m *pkg.MarginData) float64 { return m.ShortSellAmount }),
					SecurityLendRepay:  marginValue(margin, func(m *pkg.MarginData) float64 { return m.ShortRepayAmount }),
				}
				fmt.Printf("price info: %+v\n", price)
				time.Sleep(20 * time.Second)
				prices = append(prices, price)

			}

			if len(prices) > 0 {
				if err := s.stockPriceDao.BatchUpsert(ctxWithoutCancel, prices); err != nil {
					fmt.Println(fn, "batch upsert err:", err)
					return
				}
			}
			page++
			prices = prices[:0] // Clear the slice for the next iteration
		}
	}()

	return &v1.UpdateStockPriceReply{}, nil
}

func compactTradeDate(tradeDate string) string {
	return strings.ReplaceAll(tradeDate, "-", "")
}

func marginBalanceValue(margin *pkg.MarginData) float64 {
	if margin == nil {
		return 0
	}
	if margin.TotalMarginBalance != 0 {
		return margin.TotalMarginBalance
	}
	if margin.ShortBalance != 0 {
		return margin.MarginBalance + margin.ShortBalance
	}
	return margin.MarginBalance
}

func marginValue(margin *pkg.MarginData, getter func(*pkg.MarginData) float64) float64 {
	if margin == nil {
		return 0
	}
	return getter(margin)
}
