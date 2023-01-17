package utils

import (
	"errors"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

const (
	broker  = "158.160.19.212:9092"
	group   = "bmadzhuga-consumer-group"
	timeout = 6000
)

type Kafka struct {
	Topic    string
	Producer *kafka.Producer
	Consumer *kafka.Consumer
}

func (client *Kafka) Connect() error {
	kafkaProducer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": broker, "acks": 1})

	if err != nil {
		return err
	}

	client.Producer = kafkaProducer

	kafkaConsumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  broker,
		"group.id":           group,
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": false,
	})

	if err != nil {
		return err
	}

	client.Consumer = kafkaConsumer

	return nil
}

func (client *Kafka) Send(longUrl string, tinyUrl string) error {

	var err = client.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &client.Topic, Partition: kafka.PartitionAny},
		Value:          []byte(longUrl + ":" + tinyUrl)},
		nil,
	)

	if err != nil {
		errors.New("Can't send message")
	}

	fmt.Println("Message send to topic: %v", client.Topic)

	return nil
}

func (client *Kafka) Enable() error {
	if client.Consumer == nil {
		return errors.New("Empty consumer")
	}

	err := client.Consumer.SubscribeTopics([]string{client.Topic}, nil)

	if err != nil {
		return err
	}

	fmt.Println("Start reading from topic: %v", client.Topic)

	for {
		msg, err := client.Consumer.ReadMessage(-1)
		if err == nil {
			fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
		} else {
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
			break
		}
	}

	return errors.New("Error while reading")
}
