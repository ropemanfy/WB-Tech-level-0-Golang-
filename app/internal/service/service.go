package service

import (
	"L0/app/internal/config"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

type service struct {
	app   *fiber.App
	cache Cache
}

func NewService(app *fiber.App, cache Cache) service {
	return service{app: app, cache: cache}
}

func (svc *service) Start(cfg *config.App) {
	svc.app.Get("/models/:id", func(ctx *fiber.Ctx) error {
		v, err := svc.cache.Get(ctx.Params("id"))
		if err != nil {
			return err
		}
		return ctx.JSON(v)
	})
	port := fmt.Sprintf(":%s", cfg.Port)
	log.Fatal(svc.app.Listen(port))
}
