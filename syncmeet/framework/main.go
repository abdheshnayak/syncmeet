package framework

import (
	"context"
	"fmt"

	"github.com/abdheshnayak/syncmeet/syncmeet/app"
	"github.com/abdheshnayak/syncmeet/syncmeet/domain"
	"github.com/abdheshnayak/syncmeet/syncmeet/env"
	"github.com/gofiber/fiber/v2"
	"github.com/kloudlite/api/pkg/repos"
	"go.uber.org/fx"
)

type fm struct {
	*env.Env
}

func (f *fm) GetMongoConfig() (url, dbName string) {
	return f.MongoDbUri, f.MongoDbDatabse
}

var Module = fx.Module("framework",
	fx.Provide(func(ev *env.Env) *fm {
		return &fm{Env: ev}
	}),

	fx.Provide(func() *env.Env {
		return env.GetEnvOrDie()
	}),

	repos.NewMongoClientFx[*fm](),

	fx.Provide(func() *fiber.App {
		app := fiber.New()
		return app
	}),

	fx.Invoke(func(lf fx.Lifecycle, app *fiber.App, env *env.Env) {
		lf.Append(fx.Hook{
			OnStart: func(context.Context) error {
				return app.Listen(fmt.Sprintf(":%d", env.Port))
			},
			OnStop: func(context.Context) error {
				return app.Shutdown()
			},
		})
	}),

	domain.Module,
	app.Module,
)
