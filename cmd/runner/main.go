package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/paniccaaa/runner/internal/app"
	"github.com/paniccaaa/runner/internal/config"
)

func main() {
	cfg := config.MustLoad()

	// init logger
	log := setupLogger(cfg.Env)
	log.Info("starting app", slog.Int("cfg", cfg.GRPC.Port))

	// init app
	app := app.NewApp(log, cfg.GRPC.Port, cfg.DbURI)

	// start grpc-server
	go app.GRPCServer.MustRun()

	// testOutput, err := execute.ExecuteCode("package main\nimport \"fmt\"\nfunc main() {\n    fmt.Println(\"hello wrld\")\n}")
	// if err != "" {
	// 	log.Error("failed to exec code", slog.String("err", err))
	// }

	// fmt.Println(testOutput)

	// gracefull shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	app.GRPCServer.Stop()

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
