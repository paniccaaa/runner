package runner

import (
	"context"
	"errors"
	"log/slog"

	"github.com/paniccaaa/runner/internal/domain/models"

	ssoGrpc "github.com/paniccaaa/runner/internal/clients/sso/grpc"
	"github.com/paniccaaa/runner/internal/lib/execute"
	extracterr "github.com/paniccaaa/runner/internal/lib/extract-err"
)

type RunnerService struct {
	log       *slog.Logger
	storage   Storage
	ssoClient *ssoGrpc.Client
}

type Storage interface {
	GetCodeByID(ctx context.Context, id int64) (models.SharedCode, error)
	SaveCode(ctx context.Context, code, output, extractedError string) (int64, error)
	DeleteCode(ctx context.Context, id int64) error
}

func NewRunnerService(log *slog.Logger, storage Storage, ssoClient *ssoGrpc.Client) *RunnerService {
	return &RunnerService{
		log:       log,
		storage:   storage,
		ssoClient: ssoClient,
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

func (s *RunnerService) CheckAdmin(ctx context.Context, id, userID int64) (bool, error) {
	isAdmin, err := s.ssoClient.IsAdmin(ctx, userID)
	if err != nil {
		return false, err
	}

	if isAdmin {
		s.log.Info("user is admin, trying to delete code...")

		err := s.storage.DeleteCode(ctx, id)
		if err != nil {
			s.log.Error("failed to delete code", slog.String("error", err.Error()))

			return false, err
		}

		return true, nil
	}

	s.log.Warn("user is not admin", slog.Int64("user_id", userID))

	return false, nil
}
