package migrations

import (
	"database/sql"
	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upCreateTable, downCreateTable)
}

func upCreateTable(tx *sql.Tx) error {
	_, err := tx.Exec(
		`CREATE TABLE IF NOT EXISTS links (
    				id bigserial CONSTRAINT mylinks_pk PRIMARY KEY,
    				created_at TIMESTAMP,
    				hashed VARCHAR,
    				original VARCHAR
		);
	`)
	if err != nil {
		return err
	}
	return nil
}

func downCreateTable(tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	_, err := tx.Exec(`DROP TABLE IF EXISTS links;`)
	if err != nil {
		return err
	}
	return nil
}
