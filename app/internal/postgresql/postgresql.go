package postgresql

import (
	"L0/app/internal/config"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
)

type Postgres interface {
	GetClient() (conn *pgx.Conn, err error)
	Start(ctx context.Context)
	Shutdown(ctx context.Context)
}

type storageConfig struct {
	url  string
	conn *pgx.Conn
}

func NewClient(cfg *config.Postgresql) Postgres {
	url := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
	)
	return &storageConfig{url: url}
}

func (s *storageConfig) GetClient() (conn *pgx.Conn, err error) {
	if s.conn == nil {
		err = fmt.Errorf("missing connection")
		return
	}
	conn = s.conn
	return
}

func (s *storageConfig) Start(ctx context.Context) {
	if err := s.connect(ctx); err != nil {
		log.Println("failed to connect to database, starting the reconnection")
		s.reConnect(ctx)
	}
}

func (s *storageConfig) Shutdown(ctx context.Context) {
	s.conn.Close(ctx)
}

func (s *storageConfig) connect(ctx context.Context) error {
	conn, err := pgx.Connect(ctx, s.url)
	if err != nil {
		return err
	}
	s.conn = conn
	return nil
}

func (s *storageConfig) reConnect(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(5 * time.Second):
			if err := s.connect(ctx); err != nil {
				log.Println("error connect to database")
				continue
			}
			return
		}
	}
}
