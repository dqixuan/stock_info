package server

import (
	healthv1 "github.com/dqixuan/stock_info/api/health/v1"
	stockv1 "github.com/dqixuan/stock_info/api/stock/v1"
	"github.com/dqixuan/stock_info/configs"
	"github.com/dqixuan/stock_info/internal/service"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(
	c *configs.HttpServer,
	stockService *service.StockService,
	healthService *service.HealthService,
) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Addr != "" {
		opts = append(opts, http.Address(c.Addr))
	}
	if c.Timeout != 0 {
		opts = append(opts, http.Timeout(c.Timeout))
	}
	srv := http.NewServer(opts...)
	stockv1.RegisterStockServiceHTTPServer(srv, stockService)
	healthv1.RegisterHealthHTTPServer(srv, healthService)
	return srv
}
