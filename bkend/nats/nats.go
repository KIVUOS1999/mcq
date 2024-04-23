package nats

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/nats-io/nats.go"
)

type NatsStruct struct {
	natsConn *nats.Conn
}

func NewNats() (*NatsStruct, error) {
	natsUrl := strings.TrimSpace(os.Getenv("NATS_URL"))

	log.Println("checking for nats in ", natsUrl)

	nc, err := nats.Connect(natsUrl)

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

func (nc *NatsStruct) PlayerStartGame(roomID string, endTime int) error {
	err := nc.natsConn.Publish(roomID, []byte("startGame:"+strconv.Itoa(endTime/1000)))
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

func (nc *NatsStruct) Submission(roomID string, resultData []byte) error {
	err := nc.natsConn.Publish(roomID, resultData)
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
