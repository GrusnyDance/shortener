package repository

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type Instance struct {
	Db *pgxpool.Pool
}

type Link struct {
	Id       int
	Created  time.Time
	Hash     string
	Original string
}

func (i *Instance) ReturnLink(ctx context.Context, hashed string) (string, error) {
	var link Link
	var err error
	row := i.Db.QueryRow(ctx, "SELECT * FROM links WHERE hashed=$1 LIMIT 1;", hashed)
	if err = row.Scan(&link); err != nil {
		return "", errors.New("scan failed")
	}
	return link.Original, err
}

func (i *Instance) CheckIfHashedExists(ctx context.Context, hashed string) error {
	var link Link
	var err error
	row := i.Db.QueryRow(ctx, "SELECT * FROM links WHERE hashed=$1 LIMIT 1;", hashed)
	if err = row.Scan(&link); err == pgx.ErrNoRows {
		return errors.New("link not found")
	}
	return err
}

func (i *Instance) CreateLink(ctx context.Context, hashed string, original string) error {
	var err error
	_, err = i.Db.Exec(ctx, "INSERT INTO links (created_at, hashed, original) VALUES ($1, $2, $3);",
		time.Now(), hashed, original)
	return err
}
