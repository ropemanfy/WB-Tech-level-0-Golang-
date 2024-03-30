package main

import (
	"L0/app/internal/cache"
	"L0/app/internal/config"
	"L0/app/internal/db"
	"L0/app/internal/postgresql"
	"L0/app/internal/service"
	subhandlers "L0/app/internal/sub-handlers"
	"L0/app/internal/subscriber"
	"context"

	"github.com/gofiber/fiber/v2"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.GetCongfig()

	pg := postgresql.NewClient(&cfg.Postgresql)
	pg.Start(ctx)
	defer pg.Shutdown(ctx)

	database := db.NewDB(pg)

	cache := cache.NewCache(ctx, database)

	nat := subscriber.NewClient(&cfg.Nats)
	nat.Start(ctx)
	defer nat.Shutdown()
	natHandler := subhandlers.New(database, cache)
	nat.Subscribe(natHandler.MessageHandler(ctx))

	app := fiber.New()
	appservice := service.NewService(app, cache)
	appservice.Start(&cfg.App)
}
