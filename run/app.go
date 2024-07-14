package run

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"proxy/config"
	"proxy/internal/modules/controller"
	"proxy/internal/modules/service"
	"proxy/internal/modules/storage"

	"google.golang.org/grpc"
	"net"
)

type App struct {
	log        *zap.SugaredLogger
	gRPCServer *grpc.Server
	port       int
	config     *config.Config
}

func NewApp(log *zap.SugaredLogger, db *sqlx.DB, cfg *config.Config) *App {
	app := App{
		log:    log,
		config: cfg,
		port:   cfg.Local.Port,
	}
	gRPCServer := grpc.NewServer()
	Storage := storage.NewStorage(db)
	api := storage.NewApi()
	Service := service.NewService(Storage, api)

	controller.Register(gRPCServer, Service)
	controller.NewServerAPI(Service)
	app.gRPCServer = gRPCServer
	return &app
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "grpcapp.Run"

	log := a.log.With(
		"operation", op,
		"port", a.port,
	)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("grpc server is running", "address", l.Addr().String())
	if err = a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"
	a.log.With(
		"operation", op,
		"port", a.port,
	).Info("grpc server is stopping")

	a.gRPCServer.GracefulStop()
}
