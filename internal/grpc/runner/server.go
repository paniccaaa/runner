package runner

import (
	"context"
	"fmt"

	run "github.com/paniccaaa/protos/gen/golang/runner"
	"google.golang.org/grpc"
)

type Runner interface {
	RunCode(ctx context.Context, code string) (string, string, error)
	ShareCode(ctx context.Context, code string) (int64, error)
	GetCodeByID(ctx context.Context, id int64) (string, string, string, error)
}

type serverAPI struct {
	run.UnimplementedRunnerServer
	runner Runner
}

// TODO: implement grpc layer

func Register(gRPC *grpc.Server, runner Runner) {
	run.RegisterRunnerServer(gRPC, &serverAPI{runner: runner})
}

func (s *serverAPI) RunCode(ctx context.Context, req *run.CodeRequest) (*run.CodeResponse, error) {
	output, stderr, err := s.runner.RunCode(ctx, req.GetCode())
	if err != nil {
		return &run.CodeResponse{Code: req.GetCode(), Output: output, Error: stderr}, nil
		//status.Error(codes.Internal, "failed to run code")
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
