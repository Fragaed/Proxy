package service

import (
	"context"
	"fmt"
	"proxy/internal/models"
)

type Service struct {
	storage Storager
	API     API
}

//go:generate go run github.com/vektra/mockery/v2@v2.35.4 --name=Storager
type Storager interface {
	Add(response models.Response) error
	Healthcheck() error
}

//go:generate go run github.com/vektra/mockery/v2@v2.35.4 --name=API
type API interface {
	DoReq() (*models.Response, error)
	CheckAPI() error
}

func NewService(storage Storager, api API) *Service {
	return &Service{
		storage: storage,
		API:     api,
	}
}

func (s *Service) Get(ctx context.Context) (*models.Response, error) {
	resp, err := s.API.DoReq()
	if err != nil {
		return nil, err
	}

	if err = s.storage.Add(*resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *Service) Health(ctx context.Context) (*models.Health, error) {
	var resp models.Health
	if err := s.storage.Healthcheck(); err != nil {
		return nil, fmt.Errorf("healthcheck error: %w", err)
	}
	if err := s.API.CheckAPI(); err != nil {
		return nil, fmt.Errorf("healthcheck error: %w", err)
	}
	resp.DBStatus = "Postgres ready"
	resp.AppStatus = "App ready"
	resp.APIStatus = "API ready"

	return &resp, nil
}
