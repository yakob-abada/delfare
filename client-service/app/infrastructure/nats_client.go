package infrastructure

import (
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/yakob-abada/delfare/client-service/config"
)

func NewNATSClient(cfg *config.Config) *nats.Conn {
	nc, err := nats.Connect(
		cfg.NATSURL,
		nats.UserInfo(cfg.NATSUsername, cfg.NATSPassword),
	)

	if err != nil {
		panic(fmt.Sprintf("Failed to connect to NATS: %v", err))
	}

	return nc
}
