package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	Pool *pgxpool.Pool
	*Queries
}

func NewStore(ctx context.Context, databaseURL string) (*Store, error) {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	return &Store{
		Pool:    pool,
		Queries: New(pool),
	}, nil
}

func (s *Store) Close() {
	s.Pool.Close()
}

func (s *Store) WithTx(ctx context.Context, fn func(context.Context, Querier) error) error {
	tx, err := s.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := New(tx)

	if err := fn(ctx, qtx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

type StoreInterface interface {
	Querier
	WithTx(ctx context.Context, fn func(context.Context, Querier) error) error
}
