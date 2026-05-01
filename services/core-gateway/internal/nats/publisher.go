package nats

import (
	"encoding/json"
	"time"

	natsgo "github.com/nats-io/nats.go"
)

// Publisher wraps an optional NATS connection. When NATS_URL env is unset,
// all Publish calls are silent no-ops so the gateway starts cleanly offline.
type Publisher struct {
	nc *natsgo.Conn
}

func NewPublisher(natsURL string) (*Publisher, error) {
	if natsURL == "" {
		return &Publisher{nc: nil}, nil
	}
	nc, err := natsgo.Connect(natsURL, natsgo.MaxReconnects(5), natsgo.ReconnectWait(2*time.Second))
	if err != nil {
		return nil, err
	}
	return &Publisher{nc: nc}, nil
}

func (p *Publisher) Publish(subject string, payload any) error {
	if p == nil || p.nc == nil {
		return nil
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return p.nc.Publish(subject, data)
}

func (p *Publisher) Close() {
	if p != nil && p.nc != nil {
		p.nc.Drain()
	}
}
