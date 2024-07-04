package postgres

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/paniccaaa/runner/internal/domain/models"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(connStr string) (*Storage, error) {
	const op = "storage.postgres.NewStorage"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

// func (s *Storage) SaveCode(ctx context.Context, code, output, errOutput string) (int64, error) {
// 	panic("implement me")
// }

// TODO: implement storage layer

func (s *Storage) GetCodeByID(ctx context.Context, id int64) (models.SharedCode, error) {
	panic("implement me")
}

func (s *Storage) ShareCode(ctx context.Context, code string) (string, error) {
	panic("implement me")
}
