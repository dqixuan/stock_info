package configs

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewHttpCfg,
	NewGrpcCfg,
	NewMysqlCfg,
)
