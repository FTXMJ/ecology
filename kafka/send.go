package kafka

import (
	"context"
	"encoding/json"
	kfk "github.com/segmentio/kafka-go"
)

var Writer *kfk.Writer

func init() {
	config := kfk.WriterConfig{}
	config.Brokers = []string{"127.0.0.1:9092"}
	config.Topic = "ecology"
	config.Balancer = &kfk.LeastBytes{}
	Writer = kfk.NewWriter(config)
}

func SendMsg(data interface{}, topic string, key string) bool {
	bytes, _ := json.Marshal(data)
	err := Writer.WriteMessages(context.Background(), kfk.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: bytes,
	})
	if err != nil {
		return false
	}
	return true
}
