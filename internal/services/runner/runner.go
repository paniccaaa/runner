package runner

import (
	"context"
	"errors"
	"log/slog"

	"github.com/paniccaaa/runner/internal/domain/models"
	"github.com/paniccaaa/runner/internal/lib/execute"
	extracterr "github.com/paniccaaa/runner/internal/lib/extract-err"
)

type RunnerService struct {
	log     *slog.Logger
	storage Storage
}

type Storage interface {
	GetCodeByID(ctx context.Context, id int64) (models.SharedCode, error)
	SaveCode(ctx context.Context, code, output, extractedError string) (int64, error)
}

func NewRunnerService(log *slog.Logger, storage Storage) *RunnerService {
	return &RunnerService{
		log:     log,
		storage: storage,
	}
}

func (s *RunnerService) RunCode(ctx context.Context, code string) (string, string, error) {
	output, stderr := execute.ExecuteCode(code)
	if stderr != "" {
		extractedError := extracterr.ExtractSyntaxError(stderr)
		s.log.Error("failed to exec code", slog.String("err", extractedError))

		return output, extractedError, errors.New(stderr)
	}

	return output, "", nil
}

func (s *RunnerService) ShareCode(ctx context.Context, code string) (int64, error) {
	output, stderr := execute.ExecuteCode(code)
	extractedError := ""
	if stderr != "" {
		extractedError = extracterr.ExtractSyntaxError(stderr)
		s.log.Error("failed to exec code", slog.String("err", extractedError))

		return 0, errors.New(stderr)
	}

	id, err := s.storage.SaveCode(ctx, code, output, extractedError)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *RunnerService) GetCodeByID(ctx context.Context, id int64) (string, string, string, error) {
	sharedCode, err := s.storage.GetCodeByID(ctx, id)
	if err != nil {
		return "", "", "", err
	}

	return sharedCode.Code, sharedCode.Output, sharedCode.ErrOutput, nil
}
