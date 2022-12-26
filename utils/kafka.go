package utils

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Kafka struct {
	Topic    string
	Producer *kafka.Producer
	Consumer *kafka.Consumer
}

func (client *Kafka) Connect() error {
	kafkaProducer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "158.160.19.212:9092", "acks": 1})

	if err != nil {
		return err
	}

	client.Producer = kafkaProducer

	kafkaConsumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  "158.160.19.212:9092",
		"group.id":           "bmadzhuga-consumer-group",
		"enable.auto.commit": false,
	})

	if err != nil {
		return err
	}

	client.Consumer = kafkaConsumer

	return nil
}

func (client *Kafka) Send(longUrl string, tinyUrl string) error {
	chanel := make(chan kafka.Event, 10000)

	var err = client.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &client.Topic, Partition: kafka.PartitionAny},
		Value:          []byte(longUrl + ":" + tinyUrl)},
		chanel,
	)

	return err
}
