package main

import (
	_ "github.com/lib/pq"
	_ "github.com/prometheus/client_golang/prometheus"
	"urlShortener/utils"
	_ "urlShortener/utils"

	_ "github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	redis := utils.Redis{Cluster: "158.160.9.8:26379"}
	err := redis.Connect()

	if err != nil {
		panic(err)
	}
}
