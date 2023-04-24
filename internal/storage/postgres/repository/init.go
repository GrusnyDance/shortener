package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
	_ "shortener/internal/storage/postgres/migrations"
)

func Init() (*Instance, error) {
	// Create config
	poolConfig, err := NewPoolConfig()
	if err != nil {
		return nil, fmt.Errorf("pool config error: %v\n", err)
	}
	// Max connections waiting
	poolConfig.MaxConns = 10

	// Create pool of connections
	pool, err := NewConnection(poolConfig)
	if err != nil {
		return nil, fmt.Errorf("connect to database failed: %v\n", err)
	}

	// Check connection
	_, err = pool.Exec(context.Background(), ";")
	if err != nil {
		return nil, fmt.Errorf("ping failed: %v\n", err)
	}

	// Apply migrations with standard driver database/sql
	err = ApplyMigrations(poolConfig)
	if err != nil {
		pool.Close()
		return nil, errors.New("migrations failed to apply")
	}

	return &Instance{
		Db: pool,
	}, nil
}

// Wrapper for pool connection
func NewConnection(poolConfig *pgxpool.Config) (*pgxpool.Pool, error) {
	conn, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func ApplyMigrations(poolConfig *pgxpool.Config) error {
	mdb, err := sql.Open("postgres", poolConfig.ConnString())
	if err != nil {
		return err
	}
	err = mdb.Ping()
	if err != nil {
		return err
	}
	err = goose.Up(mdb, "./internal/storage/postgres/migrations")
	if err != nil {
		return err
	}
	mdb.Close()
	return nil
}
