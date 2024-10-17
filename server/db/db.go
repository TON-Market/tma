package db

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"sync"
	"time"
)

var (
	ErrOpenTransaction   = errors.New("open transaction failed")
	ErrTransactionFailed = errors.New("transaction failed")
	ErrCommitTransaction = errors.New("transaction commit failed")
)

var (
	once      sync.Once
	singleton *pgxpool.Pool
)

const dbUrl = "postgresql://postgres:password@localhost:5432/tma"

func Get() *pgxpool.Pool {
	once.Do(func() {
		config, err := pgxpool.ParseConfig(dbUrl)
		if err != nil {
			log.Fatalf("parse db config failed: %v", err)
		}

		config.MaxConns = 100
		config.MinConns = 20
		config.MaxConnLifetime = time.Hour
		config.MaxConnIdleTime = 5 * time.Minute
		config.HealthCheckPeriod = 1 * time.Minute

		pool, err := pgxpool.NewWithConfig(context.Background(), config)
		if err != nil {
			log.Fatalf("pool connection create failed: %v", err)
		}

		if _, err = pool.Exec(context.Background(), q); err != nil {
			log.Fatalln(err)
		}

		singleton = pool
	})

	return singleton
}
