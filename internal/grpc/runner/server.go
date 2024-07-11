package runner

import (
	"context"
	"fmt"

	run "github.com/paniccaaa/protos/gen/golang/runner"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RunnerProvider interface {
	RunCode(ctx context.Context, code string) (string, string, error)
	ShareCode(ctx context.Context, code string) (int64, error)
	GetCodeByID(ctx context.Context, id int64) (string, string, string, error)
	CheckAdmin(ctx context.Context, id, userId int64) (bool, error)
}

type serverAPI struct {
	run.UnimplementedRunnerServer
	runner RunnerProvider
}

func Register(gRPC *grpc.Server, runner RunnerProvider) {
	run.RegisterRunnerServer(gRPC, &serverAPI{runner: runner})
}

func (s *serverAPI) RunCode(ctx context.Context, req *run.CodeRequest) (*run.CodeResponse, error) {
	output, stderr, err := s.runner.RunCode(ctx, req.GetCode())
	if err != nil {
		return &run.CodeResponse{Code: req.GetCode(), Output: output, Error: stderr}, nil
	}

	return &run.CodeResponse{Code: req.GetCode(), Output: output, Error: ""}, nil
}

func (s *serverAPI) ShareCode(ctx context.Context, req *run.CodeRequest) (*run.ShareResponse, error) {
	id, err := s.runner.ShareCode(ctx, req.GetCode())
	if err != nil {
		return nil, fmt.Errorf("failed to share code: %w", err)
	}

	return &run.ShareResponse{Id: id}, nil
}

func (s *serverAPI) GetCodeByID(ctx context.Context, req *run.IdRequest) (*run.CodeResponse, error) {
	code, output, errOutput, err := s.runner.GetCodeByID(ctx, req.GetId())
	if err != nil {
		return nil, fmt.Errorf("failed to get code by id: %w", err)
	}

	return &run.CodeResponse{Code: code, Output: output, Error: errOutput}, nil
}

func (s *serverAPI) DeleteCode(ctx context.Context, req *run.DeleteCodeRequest) (*run.DeleteCodeResponse, error) {
	isAdmin, err := s.runner.CheckAdmin(ctx, req.GetId(), req.GetUserId())
	if err != nil {
		return nil, err
	}

	if !isAdmin {
		return nil, status.Error(codes.Canceled, "forbidden")
	}
	return &run.DeleteCodeResponse{Success: true}, nil
}
