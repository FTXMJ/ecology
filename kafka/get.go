package kafka

import (
	"context"
	"ecology/controllers"
	"ecology/models"
	"encoding/json"
	"fmt"
	kfk "github.com/segmentio/kafka-go"
)

var Read *kfk.Reader

func init() {
	config := kfk.ReaderConfig{
		Brokers:  []string{"127.0.0.1:9092"},
		GroupID:  "mrsf",
		Topic:    "ecology",
		MinBytes: 0,
		MaxBytes: 10e6,
		MaxWait:  0,
	}
	Read = kfk.NewReader(config)
}

func GetMsg() {
	fmt.Println("start")
	for {
		message, e := Read.ReadMessage(context.Background())
		if e != nil {
			fmt.Println(e.Error())
			continue
		}
		if string(message.Key) == "mrsf" {

			var data models.User
			json.Unmarshal(message.Value, &data)
			fmt.Println("mrsf:", data)
			controllers.DailyDividendAndReleaseTest(data)

		} else if string(message.Key) == "peer" {

			var data []models.User
			json.Unmarshal(message.Value, data)
			fmt.Println("peer:", data)
			controllers.PeerABounsHandler(data)

		}
	}
}

func AllTheTimeListen() {
	for {
		GetMsg()
	}
}
