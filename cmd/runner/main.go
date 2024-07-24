package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/paniccaaa/runner/internal/app"
	"github.com/paniccaaa/runner/internal/clients/gateway"
	ssoGrpc "github.com/paniccaaa/runner/internal/clients/sso/grpc"
	"github.com/paniccaaa/runner/internal/config"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	cfg := config.MustLoad()

	// init logger
	log := setupLogger(cfg.Env)
	log.Info("starting app", slog.String("env", cfg.Env))

	//init sso-client
	ssoClient, err := ssoGrpc.New(
		context.Background(),
		log,
		cfg.Clients.SSO.Address,
		cfg.Clients.SSO.Timeout,
		cfg.Clients.SSO.RetriesCount,
	)
	if err != nil {
		log.Error("failed to init sso client", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// init app
	app := app.NewApp(log, cfg.GRPC.Port, ssoClient, cfg.DbURI, cfg.RedisADDR)

	// start grpc-server
	go app.GRPCServer.MustRun()

	// start grpc-gateway
	gwServer := gateway.MustRunGRPCGateway(cfg.GRPC.Port, log)
	go gateway.StartGateway(gwServer, log)

	// expose prometheus metrics
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		log.Info("starting metrics srv", slog.String("port", cfg.PrometheusADDR))
		http.ListenAndServe(cfg.PrometheusADDR, nil)
	}()

	// gracefull shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	app.GRPCServer.Stop()
	gateway.StopGateway(gwServer, log)

	log.Info("App stopped")
}

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default: //if env config is invalid, set prod settings by default due to security
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
