package store

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose"
	"github.com/wycliff-ochieng/internal/config"
)

type Postgis struct {
	db *pgxpool.Pool
	//dbtx sqlc.DBTX
}

func NewPostgis(ctx context.Context, cfg *config.Config) (*Postgis, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.DB_USER,
		cfg.DB_PASSWORD,
		cfg.DB_HOST,
		cfg.DB_PORT,
		cfg.DB_NAME,
		cfg.DB_SSLMODE,
	)

	dbPool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection due to: %v", err)
	}

	db := stdlib.OpenDBFromPool(dbPool)

	if err := goose.SetDialect("postrges"); err != nil {
		log.Fatalf("error setting dialect: %s", err)
	}

	if err := goose.Up(db, "path/to/migrations"); err != nil {
		log.Fatalf("error spinning up goose: %s", err)
	}
	return &Postgis{dbPool}, nil
}

//Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
//Query(context.Context, string, ...interface{}) (pgx.Rows, error)
//QueryRow(context.Context, string, ...interface{}) pgx.Row

func (pg *Postgis) Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	return pg.db.Exec(ctx, query, args...)
}

func (pg *Postgis) Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	return pg.db.Query(ctx, query, args...)
}

func (pg *Postgis) QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	return pg.db.QueryRow(ctx, query, args...)
}

func (pg *Postgis) Init(ctx context.Context) error {
	return pg.CreateNewLoci(ctx)
}

func (pg *Postgis) CreateNewLoci(ctx context.Context) error {
	query := ``

	_, err := pg.db.Exec(ctx, query)
	if err != nil {
		return err
	}
	return err
}
