package nats

import (
	"log"
	"strings"

	"github.com/nats-io/nats.go"
)

type NatsStruct struct {
	natsConn *nats.Conn
}

func NewNats() (*NatsStruct, error) {
	nc, err := nats.Connect(nats.DefaultURL)

	if err != nil {
		return nil, err
	}

	return &NatsStruct{natsConn: nc}, nil
}

func (nc *NatsStruct) PlayerJoinMessage(roomID string, players []string) error {
	log.Println("room", roomID, "players", strings.Join(players, ";"))
	err := nc.natsConn.Publish(roomID, []byte(strings.Join(players, ";")))
	if err != nil {
		return err
	}

	return nil
}

func (nc *NatsStruct) PlayerStartGame(roomID string) error {
	err := nc.natsConn.Publish(roomID, []byte("startGame:"+roomID))
	if err != nil {
		return err
	}

	return nil
}

func (nc *NatsStruct) EndGame(roomID string) error {
	err := nc.natsConn.Publish(roomID, []byte("submit_game"))
	if err != nil {
		return err
	}

	return nil
}

func (nc *NatsStruct) DeleteTopic(roomID string) error {
	subscription, err := nc.natsConn.Subscribe(roomID, func(msg *nats.Msg) {
		// Message handler function
	})
	if err != nil {
		return err
	}

	subscription.Unsubscribe()

	return nil
}
