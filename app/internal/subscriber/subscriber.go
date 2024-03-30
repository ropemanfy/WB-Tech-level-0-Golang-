package subscriber

import (
	"L0/app/internal/config"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/stan.go"
)

type Nats interface {
	GetClient() (conn *stan.Conn, err error)
	Start(ctx context.Context)
	Shutdown()
	Subscribe(handler stan.MsgHandler) (err error)
}

type natsConfig struct {
	subject   string
	clusterID string
	clientID  string
	natsURL   string
	conn      stan.Conn
}

func NewClient(cfg *config.Nats) Nats {
	return &natsConfig{subject: cfg.Subject, clusterID: cfg.ClusterID,
		clientID: cfg.ClientID, natsURL: cfg.NatsUrl}
}

func (n *natsConfig) GetClient() (conn *stan.Conn, err error) {
	if n.conn == nil {
		err = fmt.Errorf("missing connection")
		return
	}
	return
}

func (n *natsConfig) Start(ctx context.Context) {
	if err := n.connect(); err != nil {
		log.Println("failed to connect to nats, starting the recconection")
		n.reConnect(ctx)
	}
}

func (n *natsConfig) Shutdown() {
	n.conn.Close()
}

func (n *natsConfig) Subscribe(handler stan.MsgHandler) (err error) {
	conn := n.conn
	_, err = conn.Subscribe(n.subject, handler, stan.SetManualAckMode())
	if err != nil {
		log.Println(err)
		return
	}
	return
}

func (n *natsConfig) connect() error {
	conn, err := stan.Connect(n.clusterID, n.clientID, stan.NatsURL(n.natsURL))
	if err != nil {
		return err
	}
	n.conn = conn
	return nil
}

func (n *natsConfig) reConnect(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(5 * time.Second):
			if err := n.connect(); err != nil {
				log.Println("error connect to nats")
				continue
			}
			return
		}
	}
}
