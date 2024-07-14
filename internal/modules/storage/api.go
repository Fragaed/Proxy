package storage

import (
	"encoding/json"
	"fmt"
	"net/http"
	"proxy/internal/models"
)

type API struct{}

func NewApi() *API {
	return &API{}
}

func (a *API) DoReq() (*models.Response, error) {
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

func (a *API) CheckAPI() error {
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
