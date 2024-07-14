package service

import (
	"context"
	"proxy/internal/models"
	"proxy/internal/modules/service/mocks"
	"reflect"
	"testing"

	"github.com/stretchr/testify/mock"
)

func TestService_Get(t *testing.T) {
	asks := models.Asks{Price: "89.35"}
	bids := models.Bids{Price: "89.34"}
	type fields struct {
		storage Storager
		API     API
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    *models.Response
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				ctx: context.Background(),
			},
			want: &models.Response{
				TimeStamp: 1720966914,
				Asks:      append([]models.Asks{}, asks),
				Bids:      append([]models.Bids{}, bids),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := mocks.NewStorager(t)
			api := mocks.NewAPI(t)
			s := &Service{
				storage: storage,
				API:     api,
			}

			api.On("DoReq").Return(&models.Response{
				TimeStamp: 1720966914,
				Asks:      append([]models.Asks{}, asks),
				Bids:      append([]models.Bids{}, bids),
			}, nil)

			storage.On("Add", mock.AnythingOfType("models.Response")).Return(nil)

			got, err := s.Get(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}

			storage.AssertExpectations(t)
			api.AssertExpectations(t)
		})
	}
}

func TestService_Health(t *testing.T) {
	type fields struct {
		storage Storager
		API     API
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.Health
		wantErr bool
	}{
		{
			name: "test",
			fields: fields{
				storage: mocks.NewStorager(t),
				API:     mocks.NewAPI(t),
			},
			args: args{
				ctx: context.Background(),
			},
			want: &models.Health{
				DBStatus:  "Postgres ready",
				AppStatus: "App ready",
				APIStatus: "API ready",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := tt.fields.storage.(*mocks.Storager)
			api := tt.fields.API.(*mocks.API)
			s := &Service{
				storage: storage,
				API:     api,
			}

			storage.On("Healthcheck").Return(nil)
			api.On("CheckAPI").Return(nil)

			got, err := s.Health(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Health() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Health() got = %v, want %v", got, tt.want)
			}

			storage.AssertExpectations(t)
			api.AssertExpectations(t)
		})
	}
}
