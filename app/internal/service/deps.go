package service

import "L0/app/internal/models"

type Cache interface {
	Get(id string) (m models.Model, err error)
}
