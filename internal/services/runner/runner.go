package runner

import (
	"context"
	"log/slog"

	"github.com/paniccaaa/runner/internal/domain/models"
)

type RunnerService struct {
	log     *slog.Logger
	storage Storage
}

type Storage interface {
	// SaveCode(ctx context.Context, code, output, errOutput string) (int64, error)
	GetCodeByID(ctx context.Context, id int64) (models.SharedCode, error)
	ShareCode(ctx context.Context, code string) (string, error)
}

func NewRunnerService(log *slog.Logger, storage Storage) *RunnerService {
	return &RunnerService{
		log:     log,
		storage: storage,
	}
}

// TODO: implement service layer

func (s *RunnerService) RunCode(ctx context.Context, code string) (string, string, error) {
	panic("implement me")
}

func (s *RunnerService) ShareCode(ctx context.Context, code string) (string, error) {
	panic("implement me")
}

func (s *RunnerService) GetCodeByID(ctx context.Context, id int64) (string, string, string, error) {
	panic("implement me")
}
