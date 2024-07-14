package migrate

import (
	"database/sql"
	"fmt"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upCreateTables, downCreateTables)
}

func upCreateTables(tx *sql.Tx) error {
	// Create order_book table
	_, err := tx.Exec(`
        CREATE TABLE IF NOT EXISTS order_book (
            id SERIAL PRIMARY KEY,
            time_stamp VARCHAR(255) NOT NULL,
            ask_price VARCHAR(255) NOT NULL,
            bid_price VARCHAR(255) NOT NULL
        );
    `)
	if err != nil {
		return fmt.Errorf("could not create order_book table: %v", err)
	}

	return nil
}

func downCreateTables(tx *sql.Tx) error {
	// Drop tables in reverse order of creation to avoid foreign key constraints
	_, err := tx.Exec(`DROP TABLE IF EXISTS order_book;`)
	if err != nil {
		return fmt.Errorf("could not drop order_book table: %v", err)
	}
	return nil
}
