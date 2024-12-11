package application

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/IskanderSh/tages-task/config"
	delivery "github.com/IskanderSh/tages-task/internal/delivery/grpc"
	"github.com/IskanderSh/tages-task/internal/service"
	"github.com/IskanderSh/tages-task/internal/storage"
	pb "github.com/IskanderSh/tages-task/proto"
	"google.golang.org/grpc"
)

type Application struct {
	cfg    *config.Config
	server *grpc.Server
	pb.UnimplementedFileProviderServer
}

func NewApplication(log *slog.Logger, cfg *config.Config) (*Application, error) {
	// storages
	fileStorage := storage.NewFileStorage(log, cfg.FileStorage.DirectoryName)
	metaStorage, err := storage.NewMetaStorage(log, cfg.MetaStorage)
	if err != nil {
		return nil, err
	}

	// services
	srv := service.NewService(log, fileStorage, metaStorage)

	// handlers
	handler := delivery.NewHandler(log, srv)

	// server
	grpcServer := grpc.NewServer()
	pb.RegisterFileProviderServer(grpcServer, handler)

	return &Application{
		cfg:    cfg,
		server: grpcServer,
	}, nil
}

func (a *Application) Run() error {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", a.cfg.Application.Host, a.cfg.Application.Port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	return a.server.Serve(listener)
}

func (a *Application) Shutdown() {
	a.server.GracefulStop()
}
