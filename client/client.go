package main

import (
	"fmt"
	_ "github.com/lib/pq"
	_ "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"net/http"
	"urlShortener/utils"
	_ "urlShortener/utils"

	_ "github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {

	b, err := ioutil.ReadFile("pass.conf")
	if err != nil {
		panic(err)
	}

	key, err := ioutil.ReadFile("//users//bogdan//.ssh//id_ed25519")
	if err != nil {
		panic(err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		panic(err)
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
		panic(err)
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
		panic(err)
	}

	kafkaClient := &utils.Kafka{
		Topic: "bmadzhuga-client-a",
		Type:  "client",
	}

	err = kafkaClient.Connect()

	if err != nil {
		panic(err)
	}

	utils.Client = kafkaClient

	redis := utils.Redis{Cluster: "158.160.9.8"}
	err = redis.Connect()

	if err != nil {
		panic(err)
	}

	defer redis.Close()
	defer client.Close()
	defer kafkaClient.Close()

	utils.InitData(client)

	utils.RegPrometheus()
	http.Handle("/metrics", promhttp.Handler())

	http.HandleFunc("/", utils.HandleGet)
	http.HandleFunc("/ping", utils.HandlePing)

	fmt.Println("Server started")

	if err := http.ListenAndServe(":8081", nil); err != nil {
		panic(err)
	}

}
