package cache

import (
	"L0/app/internal/models"
	"context"
)

type Database interface {
	GetAll(ctx context.Context) ([]models.Model, error)
}
