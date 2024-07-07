package grpcapp

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/paniccaaa/runner/internal/grpc/runner"
	"github.com/paniccaaa/runner/internal/storage/postgres"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
	storage    *postgres.Storage
}

func NewApp(log *slog.Logger, port int, runnerService runner.Runner, storage *postgres.Storage) *App {
	gRPCServer := grpc.NewServer()

	// register a new service
	runner.Register(gRPCServer, runnerService)
	reflection.Register(gRPCServer)

	return &App{
		log:        log,
		port:       port,
		gRPCServer: gRPCServer,
		storage:    storage,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "grpcapp.Run"

	log := a.log.With(
		slog.String("op", op),
		slog.Int("port", a.port),
	)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("gRPC server is running", slog.String("addr", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.log.With(slog.String("op", op)).Info("stopping gRPC server", slog.Int("port", a.port))

	a.gRPCServer.GracefulStop()

	if err := a.storage.Close(); err != nil {
		a.log.Error("failed to close storage", slog.String("error", err.Error()))
	}
}
