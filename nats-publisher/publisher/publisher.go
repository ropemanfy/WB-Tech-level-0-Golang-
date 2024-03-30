package publisher

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/nats-io/stan.go"
)

const (
	clusterID = "test-cluster"
	clientID  = "test-publisher"
	subject   = "test-subject"
	natsUrl   = "nats://localhost:4222"
)

type Publisher interface {
	StartPublic(num int)
}

type pubConfig struct {
	clusterID string
	clientID  string
	subject   string
	natsURL   string
}

func NewPublisher() Publisher {
	return &pubConfig{clusterID: clusterID,
		clientID: clientID,
		subject:  subject}
}

func (p *pubConfig) public(data []byte) {
	nc, err := stan.Connect(p.clusterID, p.clientID, stan.NatsURL(p.natsURL))
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	err = nc.Publish(p.subject, data)
	if err != nil {
		log.Println(err)
	}

}

func (p *pubConfig) StartPublic(num int) {
	for i := 1; i <= num; i++ {
		path := fmt.Sprintf("testModels/%v.json", i)
		data, err := os.ReadFile(path)
		if err != nil {
			log.Println(err)
		}
		p.public(data)
		time.Sleep(2 * time.Second)
	}
}
