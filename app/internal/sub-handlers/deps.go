package subhandlers

import (
	"L0/app/internal/models"
	"context"
)

type Cache interface {
	Create(ctx context.Context, model models.Model) (err error)
}

type Database interface {
	Create(ctx context.Context, m models.Model) error
}
