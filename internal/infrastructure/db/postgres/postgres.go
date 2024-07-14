package postgres

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
	"log"
	"proxy/config"
	_ "proxy/internal/infrastructure/db/migrate"
	"time"
)

func NewPostgresDB(cfg *config.Config) (*sqlx.DB, error) {
	var dsn string
	var err error
	var dbRaw *sql.DB

	dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.Username, cfg.DB.Password, cfg.DB.DBName, cfg.DB.SSlMode)
	fmt.Println("Connecting with DSN:", dsn)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	timeoutExceeded := time.After(time.Second * cfg.DB.TimeOut)

	for {
		select {
		case <-timeoutExceeded:
			return nil, fmt.Errorf("db connection failed after %d timeout %s", 5, err)
		case <-ticker.C:
			dbRaw, err = sql.Open(cfg.DB.Driver, dsn)
			if err != nil {
				return nil, fmt.Errorf("failed to connect to db %w", err)
			}
			err = dbRaw.Ping()
			if err == nil {

				db := sqlx.NewDb(dbRaw, cfg.DB.Driver)

				err = goose.Up(dbRaw, "./")
				if err != nil {
					log.Fatal("Goose up failed ", err)
				}
				return db, nil
			}

			log.Fatal("failed to connect to the database", err)
		}
	}
}
