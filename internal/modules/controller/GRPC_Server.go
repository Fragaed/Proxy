package controller

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"log/slog"
	"proxy/internal/models"
	protobuff "proxy/pkg/proto/proxy/gen/go"
	"strconv"
	"time"
)

type ProxyService interface {
	Get(ctx context.Context) (*models.Response, error)
	Health(ctx context.Context) (*models.Health, error)
}

type ServerAPI struct {
	service ProxyService
	protobuff.ProxyServer
}

func Register(gRPC *grpc.Server, proxy ProxyService) {
	protobuff.RegisterProxyServer(gRPC, &ServerAPI{service: proxy})
}

func NewServerAPI(proxy ProxyService) *ServerAPI {
	return &ServerAPI{
		service: proxy,
	}
}

func (s *ServerAPI) GetRates(ctx context.Context, _ *emptypb.Empty) (*protobuff.GetRatesResponse, error) {
	answer, err := s.service.Get(ctx)
	if err != nil {
		slog.Info("GetRates", "error", err)
		return nil, err
	}
	timestamp := strconv.Itoa(int(answer.TimeStamp))
	resp := &protobuff.GetRatesResponse{
		TimeStamp: timestamp,
		AskPrice:  answer.Asks[0].Price,
		BidPrice:  answer.Bids[0].Price,
	}
	slog.Info("GetRates", "time", timestamp)
	return resp, nil
}

func (s *ServerAPI) Healthcheck(ctx context.Context, _ *emptypb.Empty) (*protobuff.HealthcheckResponse, error) {
	answer, err := s.service.Health(ctx)
	if err != nil {
		slog.Info("Healthcheck", "error", err)
		return nil, err
	}
	resp := &protobuff.HealthcheckResponse{
		AppStatus:      answer.AppStatus,
		PostgersStatus: answer.DBStatus,
		ApiStatus:      answer.APIStatus,
	}
	slog.Info("Healthcheck", "time", time.Now().UTC().Format(time.RFC3339Nano))
	return resp, nil
}
