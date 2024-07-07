package gateway

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	run "github.com/paniccaaa/protos/gen/golang/runner"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func MustRunGRPCGateway(gRPCPort int, log *slog.Logger) *http.Server {
	gwServer, err := RunGRPCGateway(gRPCPort, log)
	if err != nil {
		panic(err)
	}
	return gwServer
}

func RunGRPCGateway(gRPCPort int, log *slog.Logger) (*http.Server, error) {
	const op = "grpc.gateway.RunGRPCGateway"

	address := fmt.Sprintf("0.0.0.0:%d", gRPCPort)
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	mux := runtime.NewServeMux()

	// register a runner handler
	err = run.RegisterRunnerHandler(context.Background(), mux, conn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	gwServer := &http.Server{
		Addr:    ":8090",
		Handler: mux,
	}

	log.Info("gRPC-Gateway is running", slog.String("addr", ":8090"))

	return gwServer, nil
}

func StartGateway(gwServer *http.Server, log *slog.Logger) {
	if err := gwServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Error("failed to start grpc-gateway", slog.String("error", err.Error()))
		return
	}
	log.Info("gRPC-Gateway is running", slog.String("addr", ":8090"))
}

func StopGateway(gwServer *http.Server, log *slog.Logger) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := gwServer.Shutdown(ctx); err != nil {
		log.Error("failed to gracefully shutdown gRPC-Gateway", slog.String("error", err.Error()))
	}
}
