package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
	_ "shortener/internal/storage/postgres/migrations"
)

func InitRep() (*pgxpool.Pool, error) {
	// Create config
	poolConfig, err := NewPoolConfig()
	if err != nil {
		return nil, fmt.Errorf("Pool config error: %v\n", err)
	}
	// Max connections waiting
	poolConfig.MaxConns = 10

	// Create pool of connections
	pool, err := NewConnection(poolConfig)
	if err != nil {
		return nil, fmt.Errorf("Connect to database failed: %v\n", err)
	}

	// Check connection
	_, err = pool.Exec(context.Background(), ";")
	if err != nil {
		return nil, fmt.Errorf("Ping failed: %v\n", err)
	}

	// Apply migrations with standard driver database/sql
	mdb, err := sql.Open("postgres", poolConfig.ConnString())
	err = mdb.Ping()
	err = goose.Up(mdb, "./internal/storage/postgres/migrations")
	if err != nil {
		panic(err)
	}
	mdb.Close()

	return pool, nil
}

// Wrapper for pool connection
func NewConnection(poolConfig *pgxpool.Config) (*pgxpool.Pool, error) {
	conn, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
