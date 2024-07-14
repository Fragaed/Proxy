package storage

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"proxy/internal/models"
	"strconv"
)

type PostgresAdapter struct {
	db *sqlx.DB
}

func NewStorage(db *sqlx.DB) *PostgresAdapter {
	return &PostgresAdapter{
		db: db,
	}
}

func (p *PostgresAdapter) Add(response models.Response) error {
	timestamp := strconv.Itoa(int(response.TimeStamp))
	query := `INSERT INTO order_book (time_stamp, ask_price, bid_price) VALUES ($1, $2, $3)`
	_, err := p.db.Exec(query, timestamp, response.Asks[0].Price, response.Bids[0].Price)
	if err != nil {
		return fmt.Errorf("insert to DB: %w", err)
	}
	return nil
}

func (p *PostgresAdapter) Healthcheck() error {
	return p.db.Ping()
}
