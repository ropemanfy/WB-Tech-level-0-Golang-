package subhandlers

import (
	"L0/app/internal/models"
	"context"
	"encoding/json"
	"log"

	"github.com/nats-io/stan.go"
)

type Handler struct {
	db    Database
	cache Cache
}

func New(db Database, cache Cache) *Handler {
	return &Handler{db: db, cache: cache}
}

func (h *Handler) MessageHandler(ctx context.Context) stan.MsgHandler {
	return func(msg *stan.Msg) {
		var model models.Model
		insert := func() error {
			err := json.Unmarshal(msg.Data, &model)
			if err != nil {
				return err
			}
			if err = model.Validate(); err != nil {
				return err
			}
			err = h.db.Create(ctx, model)
			if err != nil {
				return err
			}
			err = h.cache.Create(ctx, model)
			if err != nil {
				return err
			}
			return err
		}

		err := insert()
		if err != nil {
			log.Println(err)
			return
		}
		if err = msg.Ack(); err != nil {
			log.Println(err)
			return
		}
	}
}
