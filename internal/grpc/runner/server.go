package runner

import (
	"context"

	run "github.com/paniccaaa/protos/gen/golang/runner"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Runner interface {
	RunCode(ctx context.Context, code string) (string, string, error)
	ShareCode(ctx context.Context, code string) (string, error)
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
	code, output, err := s.runner.RunCode(ctx, req.GetCode())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to run code")
	}
	return &run.CodeResponse{Code: code, Output: output, Error: ""}, nil
}

func (s *serverAPI) ShareCode(ctx context.Context, req *run.CodeRequest) (*run.ShareResponse, error) {
	panic("implement me")
}

func (s *serverAPI) GetCodeByID(ctx context.Context, req *run.IdRequest) (*run.CodeResponse, error) {
	panic("implement me")
}
