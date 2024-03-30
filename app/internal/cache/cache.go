package cache

import (
	"L0/app/internal/models"
	"context"
	"fmt"
	"sync"
)

type Cache interface {
	Create(ctx context.Context, model models.Model) (err error)
	Get(id string) (m models.Model, err error)
}

type cache struct {
	data map[string]models.Model
	db   Database
	sync.RWMutex
}

func NewCache(ctx context.Context, db Database) Cache {
	cache := cache{data: make(map[string]models.Model), db: db}
	cache.recovery(ctx)
	return &cache
}

func (c *cache) recovery(ctx context.Context) error {
	c.Lock()
	defer c.Unlock()
	models, err := c.db.GetAll(ctx)
	if err != nil {
		return err
	}
	for _, v := range models {
		c.data[v.OrderUid] = v
	}
	return nil
}

func (c *cache) Create(ctx context.Context, model models.Model) (err error) {
	c.Lock()
	defer c.Unlock()
	id := model.OrderUid
	fmt.Println(id)
	_, ok := c.data[id]
	if ok {
		err = fmt.Errorf("model already exists")
		return err
	}
	c.data[id] = model
	return nil
}

func (c *cache) Get(id string) (models.Model, error) {
	c.RLock()
	defer c.RUnlock()
	model, ok := c.data[id]
	if !ok {
		err := fmt.Errorf("id not found")
		return model, err
	}
	return model, nil
}
