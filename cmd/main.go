package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/IskanderSh/tages-task/config"
	"github.com/IskanderSh/tages-task/internal/application"
)

func main() {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	cfg := config.MustLoad()
	log.Info("config loaded successfully")

	app, err := application.NewApplication(log, cfg)
	if err != nil {
		panic(err)
	}

	go func() {
		log.Info("running application")
		if err := app.Run(); err != nil {
			panic(err)
		}
	}()

	// graceful shutdown
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	<-exit

	log.Info("graceful stop")
	app.Shutdown()
}
