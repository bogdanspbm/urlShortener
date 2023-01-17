package main

import (
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	_ "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"urlShortener/utils"
	_ "urlShortener/utils"

	_ "github.com/confluentinc/confluent-kafka-go/kafka"
)

var clients = map[string]bool{"bmadzhuga-client-a": true, "bmadzhuga-client-b": true}

func main() {

	b, err := ioutil.ReadFile("pass.conf")
	if err != nil {
		fmt.Print(err)
		return
	}

	key, err := ioutil.ReadFile("//users//bogdan//.ssh//id_ed25519")
	if err != nil {
		log.Fatalf("Unable to read private key: %v", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("Unable to parse private key: %v", err)
	}

	// convert bytes to string
	pass := string(b)

	server := &utils.SSH{
		Ip:     "158.160.9.8",
		User:   "bmadzhuga",
		Port:   22,
		Cert:   pass,
		Signer: signer,
	}

	err = server.Connect(utils.CERT_PUBLIC_KEY_FILE)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer server.Close()

	utils.InitConection(*server)

	client := &utils.DBConnect{
		Ip:   "localhost",
		User: "postgres",
		Name: "bmadzhuga",
		Cert: pass}

	err = client.Open()

	if err != nil {
		fmt.Println(err)
		return
	}

	kafkaClient := &utils.Kafka{
		Topic: "bmadzhuga-events",
		Type:  "master",
	}

	err = kafkaClient.Connect()

	if err != nil {
		fmt.Println(err)
		return
	}

	utils.Client = kafkaClient

	defer client.Close()
	defer kafkaClient.Consumer.Close()
	defer kafkaClient.Producer.Close()

	go listenTopic(kafkaClient)

	utils.InitData(client)

	utils.RegPrometheus()
	http.Handle("/metrics", promhttp.Handler())

	http.HandleFunc("/", utils.HandleGet)
	http.HandleFunc("/ping", utils.HandlePing)
	http.HandleFunc("/create", utils.HandlePut)

	fmt.Println("Server started")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}

}

func listenTopic(client *utils.Kafka) {

	if client.Consumer == nil {
		panic(errors.New("Empty consumer"))
	}

	fmt.Println("Start reading from topic: ", client.Topic)

	for {
		msg, err := client.Consumer.ReadMessage(-1)
		if err == nil {
			request := strings.Split(string(msg.Value), ":")
			key := request[1]
			topic := request[0]

			if !clients[topic] {
				continue
			}

			url, ok := utils.GetURLFromKey(key)

			kafkaClient := &utils.Kafka{
				Topic: topic,
				Type:  "client",
			}

			err = kafkaClient.Connect()

			if err != nil {
				panic(err)
			}

			err = kafkaClient.Send(key, url, ok)

			fmt.Printf("Success sent answer to %v as %v:%v", topic, key, url)

			if err != nil {
				panic(err)
			}
		} else {
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
			break
		}
	}

}
