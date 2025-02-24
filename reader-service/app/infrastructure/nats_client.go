package infrastructure

import (
	"fmt"

	"github.com/nats-io/nats.go"
)

func NewNATSClient(url, username, password string) *nats.Conn {
	nc, err := nats.Connect(
		url,
		nats.UserInfo(username, password),
	)

	if err != nil {
		panic(fmt.Sprintf("Failed to connect to NATS: %v", err))
	}
	return nc
}
