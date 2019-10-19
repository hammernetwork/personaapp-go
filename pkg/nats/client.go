package nats

import (
	"encoding/json"

	"github.com/cockroachdb/errors"
	"github.com/nats-io/go-nats"
)

type Bus struct {
	client *nats.Conn
}

//nolint TODO: configure with TLS https://github.com/nats-io/go-nats#tls
func New(config Config) (*Bus, error) {
	client, err := nats.Connect(config.Addr)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to connect by addr=%s", config.Addr)
	}

	return &Bus{client: client}, nil
}

func (b *Bus) Publish(subject string, msg []byte) error {
	return b.client.Publish(subject, msg)
}

func (b *Bus) Close() {
	b.client.Close()
}

// MarshalAndPublish is a handy method to publish a json encoded msg.
func (b *Bus) MarshalAndPublish(subject string, msg interface{}) error {
	msgb, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return b.Publish(subject, msgb)
}
