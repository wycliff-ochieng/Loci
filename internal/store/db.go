package store

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pressly/goose"
	"github.com/wycliff-ochieng/internal/config"
)

type Postgis struct {
	db *pgxpool.Pool
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

	db, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection due to: %v", err)
	}

	if err := goose.SetDialect(); err != nil {
		log.Fatalf("error setting dialect: %s", err)
	}

	if err := goose.Up(db, "path/to/migrations"); err != nil {
		log.Fatalf("error spinning up goose: %s", err)
	}
	return &Postgis{db}, nil
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
