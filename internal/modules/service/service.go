package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"proxy/internal/models"
)

type Service struct {
	storage Storager
}

type Storager interface {
	Add(response models.Response) error
	Healthcheck() error
}

func NewService(storage Storager) *Service {
	return &Service{
		storage: storage,
	}
}

func (s *Service) Get(ctx context.Context) (*models.Response, error) {
	resp, err := doReq()
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
	if err := chekAPI(); err != nil {
		return nil, fmt.Errorf("healthcheck error: %w", err)
	}
	resp.DBStatus = "Postgres ready"
	resp.AppStatus = "App ready"
	resp.APIStatus = "API ready"

	return &resp, nil
}

func doReq() (*models.Response, error) {
	var resp models.Response
	url := "https://garantex.org/api/v2/depth?market=usdtrub"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	req.Header.Add("Cookie", "__ddg1_=SddCSyqTtTuXc6PZS6ka")

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer res.Body.Close()

	if err = json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return nil, fmt.Errorf("decode response body: %w", err)
	}
	return &resp, nil
}

func chekAPI() error {
	url := "https://garantex.org/api/v2/depth?market=usdtrub"
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	req.Header.Add("Cookie", "__ddg1_=SddCSyqTtTuXc6PZS6ka")

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("making request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned non-200 status: %d", res.StatusCode)
	}

	return nil
}
