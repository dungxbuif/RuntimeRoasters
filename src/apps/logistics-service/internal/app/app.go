package app

import (
	"context"

	svcconfig "github.com/dungxbuif/RuntimeRoasters/apps/logistics-service/config"
	"github.com/dungxbuif/RuntimeRoasters/pkg/base"
	"github.com/dungxbuif/RuntimeRoasters/pkg/database"
	"github.com/redis/go-redis/v9"
)

type App struct {
	Base *base.App
	Cfg  *svcconfig.Config
	DB   *database.DB
	VDB  *redis.Client
}

func NewApp(
	baseApp *base.App,
	cfg *svcconfig.Config,
	db *database.DB,
	vdb *redis.Client,
) *App {
	return &App{
		Base: baseApp,
		Cfg:  cfg,
		DB:   db,
		VDB:  vdb,
	}
}

func (a *App) Run() {
	a.Base.FinalizeRoutes()

	// Register readiness check
	a.Base.RegisterReadiness(func() error {
		if err := a.DB.Ping(context.Background()); err != nil {
			return err
		}
		return nil
	})

	if a.Cfg.GRPCPort == 0 {
		a.Cfg.GRPCPort = 50055 // Using 50055 for logistics
	}

	a.Base.Run(a.Cfg.AppPort, a.Cfg.GRPCPort)
}
