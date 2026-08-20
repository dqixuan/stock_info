package dao

import "github.com/google/wire"

var DaoProviderSet = wire.NewSet(
	NewStockDao,
	NewStockPriceDao,
	NewIndustryStockDao,
	NewIndustryDao,
)
