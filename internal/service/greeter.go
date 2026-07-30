package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	v1 "github.com/dqixuan/stock_info/api/helloworld/v1"
	"github.com/dqixuan/stock_info/configs"
	"github.com/dqixuan/stock_info/internal/biz"
	"github.com/dqixuan/stock_info/internal/dao"
	"github.com/dqixuan/stock_info/internal/model"
	"github.com/dqixuan/stock_info/pkg"
	"gorm.io/gorm"
)

// GreeterService is a greeter service.
type GreeterService struct {
	v1.UnimplementedGreeterServer

	uc       *biz.GreeterUsecase
	conf     *configs.Config
	mysql    *gorm.DB
	stockDao *dao.StockDao
}

// NewGreeterService new a greeter service.
func NewGreeterService(uc *biz.GreeterUsecase, conf *configs.Config) *GreeterService {
	return &GreeterService{uc: uc, conf: conf}
}

// SayHello implements helloworld.GreeterServer.
func (s *GreeterService) SayHello(ctx context.Context, in *v1.HelloRequest) (*v1.HelloReply, error) {
	g, err := s.uc.CreateGreeter(ctx, &biz.Greeter{Hello: in.Name})
	if err != nil {
		return nil, err
	}
	return &v1.HelloReply{Message: "Hello " + g.Hello}, nil
}

func (s *GreeterService) SaveStock(ctx context.Context, in *v1.StockInfoRequest) (*v1.StockInfoReply, error) {
	go func() {
		stocks, err := pkg.FetchAllAStocks()
		if err != nil {
			fmt.Println("err:", err)
			return
		}
		status := 1
		for _, stock := range stocks {
			if strings.Contains(stock.Name, "ST") {
				status = 2
			}
			if strings.Contains(stock.Name, "*ST") {
				status = 3
			}
			err = s.stockDao.Create(&model.Stock{
				StockID:   stock.Code,
				Name:      stock.Name,
				Market:    stock.Symbol[:2],
				Status:    int8(status),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			})
			if err != nil {
				fmt.Println("err:", err)
			}
		}

	}()
	return &v1.StockInfoReply{}, nil
}
