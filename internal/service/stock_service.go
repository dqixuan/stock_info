package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	set "github.com/deckarep/golang-set/v2"

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
	industryDao   *dao.IndustryDao
}

func NewStockService(
	data *data.Data,
	stockDao *dao.StockDao,
	stockPriceDao *dao.StockPriceDao,
	industryDao *dao.IndustryDao,
) *StockService {
	return &StockService{
		mysql:         data.DB(),
		stockDao:      stockDao,
		stockPriceDao: stockPriceDao,
		industryDao:   industryDao,
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

// UpdateStockPrice 每个交易日， 更没有记录价格信息的股票信息
// curl -X POST http://localhost:8000/api/v1/stock/price/update -H "Content-Type: application/json" -d '{}'
func (s StockService) UpdateStockPrice(ctx context.Context, request *v1.UpdateStockPriceRequest) (*v1.UpdateStockPriceReply, error) {
	fn := "UpdateStockPrice"
	go func() {
		ctxWithoutCancel := context.WithoutCancel(ctx)
		allStockIDs, err := s.stockDao.SelectAllStockIDs(ctxWithoutCancel)
		if err != nil {
			fmt.Println(fn, "SelectAllStockIDs failed err:", err)
			return
		}
		fmt.Println("fn:", fn, ", allStockIDs len:", len(allStockIDs))
		recordedStockIDs, err := s.stockPriceDao.SelectStockIDsByDate(ctxWithoutCancel, time.Now().Format(time.DateOnly))
		if err != nil {
			fmt.Println(fn, "SelectStockIDsByDate failed err:", err)
			return
		}
		fmt.Println("fn:", fn, ", recordedStockIDs len:", len(recordedStockIDs))

		setA := set.NewSet(allStockIDs...)
		setB := set.NewSet(recordedStockIDs...)

		// 2. 求差集 (在 A 中，不在 B 中)
		// 注意：Difference 方法会自动去重
		diffSet := setA.Difference(setB)

		// 3. 将 Set 转换回 []string
		diffSlice := diffSet.ToSlice()
		fmt.Println(fn, "diffSlice:", diffSlice)
		for i := 0; i < len(diffSlice); i += 100 {
			var stockIDs []string
			if i+100 <= len(diffSlice) {
				stockIDs = diffSlice[i : i+100]
			} else {
				stockIDs = diffSlice[i:]
			}
			prices := make([]*model.StockPrice, 0, len(stockIDs))

			tradeDate := time.Now().Format(time.DateOnly)
			for _, stockID := range stockIDs {
				if stockID == "" {
					continue
				}

				stockData, err := pkg.GetStockPrice(stockID)
				if err != nil {
					fmt.Println(fn, "get stock data err:", stockID, err)
					continue
				}

				//margin, err := pkg.GetStockMarginByDate(stock.StockID, compactTradeDate(tradeDate))
				//if err != nil {
				//	fmt.Println(fn, "get stock margin err:", stock.StockID, err)
				//	margin = nil
				//}

				price := &model.StockPrice{
					StockID:       stockID,
					TradeDate:     tradeDate,
					OpenPrice:     stockData.OpenPrice,
					ClosePrice:    stockData.LatestPrice,
					HighPrice:     stockData.HighPrice,
					LowPrice:      stockData.LowPrice,
					Volume:        int64(stockData.Volume),
					Amount:        stockData.Turnover,
					ChangePercent: stockData.ChangePercentage,
					//FinanceBalance:     marginValue(margin, func(m *pkg.MarginData) float64 { return m.MarginBalance }),
					//FinanceBuy:         marginValue(margin, func(m *pkg.MarginData) float64 { return m.MarginBuyAmount }),
					//FinanceRepay:       marginValue(margin, func(m *pkg.MarginData) float64 { return m.MarginRepayAmount }),
					//SecurityLendVolume: marginValue(margin, func(m *pkg.MarginData) float64 { return m.ShortSellVolume }),
					//SecurityLendSell:   marginValue(margin, func(m *pkg.MarginData) float64 { return m.ShortSellAmount }),
					//SecurityLendRepay:  marginValue(margin, func(m *pkg.MarginData) float64 { return m.ShortRepayAmount }),
				}
				fmt.Printf("price info: %+v\n", price)
				//time.Sleep(2 * time.Second)
				prices = append(prices, price)
			}

			if len(prices) > 0 {
				if err = s.stockPriceDao.BatchUpsert(ctxWithoutCancel, prices); err != nil {
					fmt.Println(fn, "batch upsert err:", err)
					return
				}
			}

			prices = prices[:0]
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

// InitStockIndustryConcept 初始化股票行业概念数据
// curl -X POST http://localhost:8000/api/v1/stock/industry_concept/init   -H "Content-Type: application/json"   -d '{}'
func (s StockService) InitStockIndustryConcept(ctx context.Context, request *v1.InitStockIndustryConceptRequest) (*v1.InitStockIndustryConceptReply, error) {
	fn := "InitStockIndustryConcept"
	ctxWithoutCancel := context.WithoutCancel(ctx)
	go func() {
		industryInfo, err := pkg.GetTHSBoards("all")
		if err != nil {
			fmt.Printf("%s GetTHSBoards failed, err:%v\n", fn, err)
			return
		}
		var industries []*model.Industry
		now := time.Now()
		for _, info := range industryInfo.Data {

			industries = append(industries, &model.Industry{
				IndustryID: info.Code,
				Name:       info.Name,
				Count:      0,
				CreatedAt:  now,
				UpdatedAt:  now,
			})
		}
		if len(industries) > 0 {
			if err = s.industryDao.BatchInsert(ctxWithoutCancel, industries); err != nil {
				fmt.Println(fn, "batch upsert err:", err)
				return
			}
		}
	}()
	return &v1.InitStockIndustryConceptReply{}, nil
}
