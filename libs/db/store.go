package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SQLStore struct {
	*Queries
	connPool *pgxpool.Pool
}

func NewSQLStore(connPool *pgxpool.Pool) *SQLStore {
	return &SQLStore{
		connPool: connPool,
		Queries:  New(connPool),
	}
}

func NewPgxPool(conn string, minConn, maxConn int32) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(conn)
	if err != nil {
		return nil, err
	}
	config.MaxConns = maxConn
	config.MinConns = minConn
	return pgxpool.NewWithConfig(context.Background(), config)
}

func (store *SQLStore) Close() {
	store.connPool.Close()
}
