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

func (s *Storage) GetCodeByID(ctx context.Context, id int64) (models.SharedCode, error) {
	const query = `
			SELECT id, code, output, error
			FROM shared_codes
			WHERE id = $1;
	`
	var sharedCode models.SharedCode
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&sharedCode.ID,
		&sharedCode.Code,
		&sharedCode.Output,
		&sharedCode.ErrOutput,
	)
	if err != nil {
		return models.SharedCode{}, err
	}

	return sharedCode, nil
}

func (s *Storage) SaveCode(ctx context.Context, code, output, extractedError string) (int64, error) {
	const query = `
			INSERT INTO shared_codes (code, output, error)
			VALUES ($1, $2, $3)
			RETURNING id;
	`
	var id int64
	err := s.db.QueryRowContext(ctx, query, code, output, extractedError).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Storage) DeleteCode(ctx context.Context, id int64) error {
	deleteQuery := "DELETE FROM shared_codes WHERE id = $1"

	result, err := s.db.ExecContext(ctx, deleteQuery, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no rows affected, code with id %d not found", id)
	}

	return nil
}
