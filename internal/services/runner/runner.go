package runner

import (
	"context"
	"errors"
	"log/slog"

	"github.com/paniccaaa/runner/internal/domain/models"
	"github.com/paniccaaa/runner/internal/lib/execute"
)

type RunnerService struct {
	log     *slog.Logger
	storage Storage
}

type Storage interface {
	// SaveCode(ctx context.Context, code, output, errOutput string) (int64, error)
	GetCodeByID(ctx context.Context, id int64) (models.SharedCode, error)
	SaveCode(ctx context.Context, code string) (string, error)
}

func NewRunnerService(log *slog.Logger, storage Storage) *RunnerService {
	return &RunnerService{
		log:     log,
		storage: storage,
	}
}

// TODO: implement service layer

func (s *RunnerService) RunCode(ctx context.Context, code string) (string, string, error) {
	output, err := execute.ExecuteCode(code)
	if err != "" {
		s.log.Error("failed to exec code", slog.String("err", err))
		return "", err, errors.New(err)
	}
	return output, "", nil
}

func (s *RunnerService) ShareCode(ctx context.Context, code string) (string, error) {
	panic("implement me")
}

func (s *RunnerService) GetCodeByID(ctx context.Context, id int64) (string, string, string, error) {
	panic("implement me")
}
