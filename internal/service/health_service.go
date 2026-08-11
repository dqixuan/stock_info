package service

import (
	"context"

	v1 "github.com/dqixuan/stock_info/api/health/v1"
)

type HealthService struct{}

func NewHealthService() *HealthService {
	return &HealthService{}
}

func (s *HealthService) Check(ctx context.Context, request *v1.HealthCheckRequest) (*v1.HealthCheckReply, error) {
	return &v1.HealthCheckReply{Result: "OK"}, nil
}

var _ v1.HealthHTTPServer = (*HealthService)(nil)